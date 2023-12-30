package server

import (
	"log"
	"net/http"

	"github.com/harmony-cloud-live/server/internal/midi"
	"github.com/harmony-cloud-live/server/internal/osc"
	"nhooyr.io/websocket"
)

type HarmonyCloudServer struct {
	serveMux http.ServeMux
	midiPlayer *midi.MidiPlayer
	oscClient *osc.OscClient

	logf func(f string, v ...interface{})

	clients map[string]*Client
	leader *Client
	
	currentIndex int
	currentBeat int
	mainSequence []Chord
}

func NewHarmonyCloudServer(midiPlayer *midi.MidiPlayer, oscClient *osc.OscClient) *HarmonyCloudServer {
	h := &HarmonyCloudServer{
		midiPlayer: midiPlayer,
		oscClient: oscClient,

		logf: log.Printf,
		clients: make(map[string]*Client),

		currentIndex: 0,
		currentBeat: 0,
		mainSequence: getDummyData(),
	}
	h.serveMux.HandleFunc("/midi", h.upgradeMidiSocket)
	h.serveMux.HandleFunc("/control", h.upgradeControlSocket)

	return h
}

func (h *HarmonyCloudServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveMux.ServeHTTP(w, r)
}

func (h *HarmonyCloudServer) addClient(userId string, username string, c *websocket.Conn) {
	h.clients[userId] = &Client{
		Conn: c,
		UserId: userId,
		Username: username,
	}
	// mutex?
	if h.leader == nil || h.leader.UserId == userId {
		h.logf("auto set leader: %v", userId)
		h.leader = h.clients[userId]
	}
}

func (h *HarmonyCloudServer) getClients() []*ClientData {
	var clients []*ClientData
	for _, c := range h.clients {
		clients = append(clients, &ClientData{UserId: c.UserId, Username: c.Username})
	}
	return clients
}