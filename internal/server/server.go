package server

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/harmony-cloud-live/server/internal/midi"
	"github.com/harmony-cloud-live/server/internal/music"
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
	apiClient *music.ApiClient

	clients map[string]*Client
	leader *Client
	cachedLeaderId string
	
	state *music.PlaybackState
}

func NewHarmonyCloudServer(midiPlayer *midi.MidiPlayer, oscClient *osc.OscClient, 
	rdb *redis.Client, apiClient *music.ApiClient) (*HarmonyCloudServer, error) {
	ctx := context.Background()

	h := &HarmonyCloudServer{
		midiPlayer: midiPlayer,
		oscClient: oscClient,
		rdb: rdb,
		apiClient: apiClient,

		logf: log.Printf,
		clients: make(map[string]*Client),
	}

	h.serveMux.HandleFunc("/midi", h.upgradeMidiSocket)
	h.serveMux.HandleFunc("/control", h.upgradeControlSocket)

	h.cachedLeaderId = h.getCachedLeaderId(ctx)
	state, err := h.getCachedPlaybackState(ctx)

	if err != nil {
		h.logf("[INFO] could not get cached playback state: %v", err)
		h.state = music.NewPlaybackState()
		h.newMainSequence(ctx, "Legacy Dances Solo 2")
	} else {
		h.state = state
	}
	return h, nil
}

func (h *HarmonyCloudServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveMux.ServeHTTP(w, r)
}
