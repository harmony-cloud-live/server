package midi

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiPlayer struct {
	send func(msg midi.Message) error
	close func()
	
	noteDelay int64
	velocity uint8

	minNote uint8
	maxNote uint8
}

func NewMidiPlayer(portName string) (*MidiPlayer, error) {
	out, err := midi.FindOutPort(portName)
	if err != nil {
		return nil, fmt.Errorf("can't find midi port: %w", err)
	}
	
	send, err := midi.SendTo(out)
	if err != nil {
		midi.CloseDriver()
		return nil, fmt.Errorf("error initializing sender: %w", err)
	}
	
	return &MidiPlayer{
		send: send,
		close: midi.CloseDriver,

		noteDelay: 0,
		velocity: 64,

		minNote: 127,
		maxNote: 0,
	}, nil
}

func (m *MidiPlayer) SetNoteDelay(val int64) int64 {
	m.noteDelay = max(0, min(200, val))
	return m.noteDelay
}

func (m *MidiPlayer) GetNoteDelay() int64 {
	return m.noteDelay
}

func (m *MidiPlayer) SetVelocity(val uint8) uint8 {
	m.velocity = max(0, min(100, val))
	m.send(midi.ControlChange(0, 11, m.velocity))
	return m.velocity
}

func (m *MidiPlayer) GetVelocity() uint8 {
	return m.velocity
}

func (m *MidiPlayer) PlayChord(c []uint8) {
	effectiveDelay := m.noteDelay
	chord := make([]uint8, len(c))
	copy(chord, c)

	if effectiveDelay == 0 {
		rand.Shuffle(len(chord), func(i, j int) {
			chord[i], chord[j] = chord[j], chord[i]
		})
		effectiveDelay = 1
	}
	
	for _, note := range chord {
		m.send(midi.NoteOn(0, note, m.velocity))
		m.minNote = min(m.minNote, note)
		m.maxNote = max(m.maxNote, note)
		time.Sleep(time.Millisecond * time.Duration(effectiveDelay))
	}
}

func (m *MidiPlayer) StopChord(chord []uint8) {
	for _, note := range chord {
		m.send(midi.NoteOff(0, note))
	}
}

func (m *MidiPlayer) StopAll() {
	for i := m.minNote; i <= m.maxNote; i++ {
		m.send(midi.NoteOff(0, i))
	}
}

func (m *MidiPlayer) Close() {
	m.close()
}
