package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/harmony-cloud-live/server/internal/midi"
	"github.com/harmony-cloud-live/server/internal/osc"
	"github.com/redis/go-redis/v9"
	"nhooyr.io/websocket"
)

type HarmonyCloudServer struct {
	serveMux http.ServeMux

	midiPlayer *midi.MidiPlayer
	oscClient *osc.OscClient
	redisClient *redis.Client

	logf func(f string, v ...interface{})

	clients map[string]*Client
	leader *Client
	cachedLeaderId string
	mu sync.RWMutex
	
	currentIndex int
	currentBeat int
	mainSequence []Chord
}

func NewHarmonyCloudServer(midiPlayer *midi.MidiPlayer, oscClient *osc.OscClient, redisClient *redis.Client) *HarmonyCloudServer {
	h := &HarmonyCloudServer{
		midiPlayer: midiPlayer,
		oscClient: oscClient,
		redisClient: redisClient,

		logf: log.Printf,
		clients: make(map[string]*Client),

		currentIndex: 0,
		currentBeat: 0,
	}
	h.serveMux.HandleFunc("/midi", h.upgradeMidiSocket)
	h.serveMux.HandleFunc("/control", h.upgradeControlSocket)

	ctx := context.Background()
	h.cachedLeaderId = h.redisClient.Get(ctx, "leaderId").Val()
	cachedMainSequence, err := h.getCachedMainSequence(ctx)
	if err == nil {
		h.mainSequence = cachedMainSequence
	} else {
		h.newMainSequence(ctx)
	}

	return h
}

func (h *HarmonyCloudServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveMux.ServeHTTP(w, r)
}

func (h *HarmonyCloudServer) newMainSequence(ctx context.Context) {
    h.currentIndex = 0
    h.currentBeat = 0
    h.mainSequence = getDummyData()
	
	err := h.cacheMainSequence(ctx)
	if err != nil {
		h.logf("error caching main sequence: %v", err)
	}
}

func (h *HarmonyCloudServer) cacheMainSequence(ctx context.Context) error {
    var chordSymbols []string
    var midiValues [][]uint8
    for _, chord := range h.mainSequence {
        chordSymbols = append(chordSymbols, chord.ChordSymbol)
        midiValues = append(midiValues, chord.MidiValues)
    }

    encChordSymbols, err := json.Marshal(chordSymbols)
    if err != nil {
        return err
    }
    err = h.redisClient.Set(ctx, "mainSequenceChordSymbols", encChordSymbols, 0).Err()
    if err != nil {
        return err
    }

    encMidiValues, err := json.Marshal(midiValues)
    if err != nil {
        return err
    }
    err = h.redisClient.Set(ctx, "mainSequenceMidiValues", encMidiValues, 0).Err()
    if err != nil {
        return err
    }
	return nil
}

func (h *HarmonyCloudServer) getCachedMainSequence(ctx context.Context) ([]Chord, error) {
    encChordSymbols, err := h.redisClient.Get(ctx, "mainSequenceChordSymbols").Bytes()
    if err != nil {
        return nil, err
    }
    var chordSymbols []string
    err = json.Unmarshal(encChordSymbols, &chordSymbols)
    if err != nil {
        return nil, err
    }

    encMidiValues, err := h.redisClient.Get(ctx, "mainSequenceMidiValues").Bytes()
    if err != nil {
        return nil, err
    }
    var midiValues [][]uint8
    err = json.Unmarshal(encMidiValues, &midiValues)
    if err != nil {
        return nil, err
    }

    var mainSequence []Chord
    for i, symbol := range chordSymbols {
        mainSequence = append(mainSequence, Chord{ChordSymbol: symbol, MidiValues: midiValues[i]})
    }

    return mainSequence, nil
}


func (h *HarmonyCloudServer) addClient(ctx context.Context, userId string, username string, c *websocket.Conn) error {
	h.clients[userId] = &Client{
		Conn: c,
		UserId: userId,
		Username: username,
	}
	if h.leader == nil {
		if h.cachedLeaderId != "" && h.clients[h.cachedLeaderId] != nil {
			h.logf("auto set leader from cache: %v", h.clients[h.cachedLeaderId])
			return h.setLeader(ctx, h.clients[h.cachedLeaderId])
		} 
	}
	if h.leader != nil && h.leader.UserId == userId {
		h.logf("reconnecting leader: %v", h.clients[userId])
		return h.setLeader(ctx, h.clients[userId])
	}
	return nil
}

func (h *HarmonyCloudServer) getClients() []*ClientData {
	var clients []*ClientData
	for _, c := range h.clients {
		clients = append(clients, &ClientData{UserId: c.UserId, Username: c.Username})
	}
	return clients
}
