package server

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/harmony-cloud-live/server/internal/data"
	"github.com/harmony-cloud-live/server/internal/midi"
	"github.com/harmony-cloud-live/server/internal/osc"
	"github.com/redis/go-redis/v9"
)

type HarmonyCloudServer struct {
	serveMux http.ServeMux
	mu sync.RWMutex
	logf func(f string, v ...interface{})

	midiPlayer *midi.MidiPlayer
	oscClient *osc.OscClient
	rdb *redis.Client
	apiClient *data.ApiClient

	clients map[string]*Client
	leader *Client
	cachedLeaderId string
	
	mainSequence *[]data.Chord
	currentIndex int
	currentBeat int
	timeSignature TimeSignature
	songName string
	key data.KeySignature
	loopStart int
	loopEnd int
}

func NewHarmonyCloudServer(midiPlayer *midi.MidiPlayer, oscClient *osc.OscClient, rdb *redis.Client, apiClient *data.ApiClient) *HarmonyCloudServer {
	ctx := context.Background()

	h := &HarmonyCloudServer{
		midiPlayer: midiPlayer,
		oscClient: oscClient,
		rdb: rdb,
		apiClient: apiClient,

		logf: log.Printf,
		clients: make(map[string]*Client),

		currentIndex: 0,
		currentBeat: 0,
		timeSignature: TimeSignature{Upper: 4, Lower: 4},
		key: data.KeySignature("F"),
		loopStart: -1,
		loopEnd: -1,
	}
	h.serveMux.HandleFunc("/midi", h.upgradeMidiSocket)
	h.serveMux.HandleFunc("/control", h.upgradeControlSocket)

	h.cachedLeaderId = h.getCachedLeaderId(ctx)
	h.mainSequence = h.getCachedMainSequence(ctx)

	if h.mainSequence == nil || len(*h.mainSequence) == 0 {
		h.newMainSequence(ctx, "Legacy Dances Solo 2")
	}
	return h
}

func (h *HarmonyCloudServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveMux.ServeHTTP(w, r)
}
