package server

import (
	"encoding/json"

	"github.com/harmony-cloud-live/server/internal/music"
	ws "nhooyr.io/websocket"
)

type Client struct {
	Conn       *ws.Conn `json:"-"` 
	UserId     string   `json:"userId"`
	Username   string   `json:"username"`
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
    GetLoop                          // 14
    NewLoop                          // 15
    GetNoteDelay                     // 16
    SetNoteDelay                     // 17
    GetVelocity                      // 18
    SetVelocity                      // 19
    ManualChordDown                  // 20
    ManualChordUp                    // 21
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
    case GetLoop:
        return "GetLoop"
    case NewLoop:
        return "NewLoop"
    case GetNoteDelay:
        return "GetNoteDelay"
    case SetNoteDelay:
        return "SetNoteDelay"
    case GetVelocity:
        return "GetVelocity"
    case SetVelocity:
        return "SetVelocity"
    case ManualChordDown:
        return "ManualChordDown"
    case ManualChordUp:
        return "ManualChordUp"
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

func (c Chord) MarshalJSON() ([]byte, error) {
    return json.Marshal(c.ChordSymbol)
}

type ControlPayload struct {
    Index int `json:"index"`
    Beat  int `json:"beat"`
    Chords []music.Chord `json:"chords"`
    TimeSignature music.TimeSignature `json:"timeSignature"`
    LeaderId string `json:"leaderId"`
    Username string `json:"username"`
    Clients []*Client `json:"clients"`
    SongTitle string `json:"songTitle"`
    LoopStart int `json:"loopStart"`
    LoopEnd int `json:"loopEnd"`
    NoteDelay int64 `json:"noteDelay"`
    Velocity uint8 `json:"velocity"`
}

type ControlEvent struct {
    Type    ControlEventType
    Payload ControlPayload
}
