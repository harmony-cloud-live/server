package server

import (
	"encoding/json"

	"nhooyr.io/websocket"
)

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
    GetTimeSignature                 // 12
    NewTimeSignature                 // 13
)

func (c ControlEventType) String() string {
    switch c {
    case GetIndex:
        return "GetIndex"
    case NewIndex:
        return "NewIndex"
    case GetBeat:
        return "GetBeat"
    case NewBeat:
        return "NewBeat"
    case GetMainSequence:
        return "GetMainSequence"
    case NewMainSequence:
        return "NewMainSequence"
    case GetSettings:
        return "GetSettings"
    case NewSettings:
        return "NewSettings"
    case GetLeader:
        return "GetLeader"
    case NewLeader:
        return "NewLeader"
    case SetUsername:
        return "SetUsername"
    case GetClients:
        return "GetClients"
    case GetTimeSignature:
        return "GetTimeSignature"
    case NewTimeSignature:
        return "NewTimeSignature"
    default:
        return ""
    }
}

type MidiEvent struct {
    Type  MidiEventType
	Index int
}

type Chord struct {
    ChordSymbol string `json:"chordSymbol"`
    MidiValues  []uint8
}

type TimeSignature struct {
    Upper uint8 `json:"upper"`
    Lower uint8 `json:"lower"`
}

func (t TimeSignature) isValid() bool {
    return t.Upper > 0 && t.Lower > 0
}

func (c Chord) MarshalJSON() ([]byte, error) {
    return json.Marshal(c.ChordSymbol)
}

type ControlPayload struct {
    Index int `json:"index"`
    Beat  int `json:"beat"`
    Chords []Chord `json:"chords"`
    TimeSignature TimeSignature `json:"timeSignature"`
    LeaderId string `json:"leaderId"`
    Username string `json:"username"`
    Clients []*ClientData `json:"clients"`
}

type ControlEvent struct {
    Type    ControlEventType
    Payload ControlPayload
}
