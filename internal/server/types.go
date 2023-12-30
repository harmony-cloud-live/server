package server

import "nhooyr.io/websocket"

type Client struct {
	Conn *websocket.Conn
	UserId string
	Username string
}

type ClientData struct {
	UserId string `json:"userId"`
	Username string `json:"username"`
}

type MidiEventType byte

const (
    ChordDown MidiEventType = iota // 0
    ChordUp                        // 1
    StopAll                        // 2
)

type ControlEventType byte

const (
    GetIndex ControlEventType = iota // 0
    NewIndex                         // 1
    GetBeat                          // 2
    NewBeat                          // 3
    GetMainSequence                  // 4
    NewMainSequence                  // 5
    GetSettings                      // 6
    NewSettings                      // 7
    GetLeader                        // 8
    NewLeader                        // 9
    SetUsername                      // 10
    GetClients                       // 11
)

func (c ControlEventType) String() string {
	return [...]string{"GetIndex", "NewIndex", "GetBeat", "NewBeat", "GetMainSequence", "NewMainSequence", "GetSettings", "NewSettings", "GetLeader", "NewLeader", "SetUsername", "GetClients"}[c]
}

type MidiEvent struct {
    Type  MidiEventType
	Index int
}

type Chord struct {
    ChordSymbol string `json:"chordSymbol"`
    MidiValues  []uint8 `json:"midiValues"`
}

// Using an interface{} as it can be of several types
type ControlPayload interface{}

type ControlEvent struct {
    Type    ControlEventType
    Payload ControlPayload
}

