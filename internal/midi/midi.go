package midi

import (
	"fmt"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiPlayer struct {
	send func(msg midi.Message) error
	close func()

	mu sync.Mutex
	currentChord []uint8
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
	}, nil
}


func (m *MidiPlayer) PlayChord(chord []uint8) {
	for _, note := range chord {
		m.send(midi.NoteOn(0, note, 64))
	}

	m.mu.Lock()
	m.currentChord = chord
	m.mu.Unlock()
}

func (m *MidiPlayer) StopCurrentChord() {
	m.mu.Lock()
	currentChord := m.currentChord
	m.currentChord = nil
	m.mu.Unlock()

	for _, note := range currentChord {
		m.send(midi.NoteOff(0, note))
	}
}

func (m *MidiPlayer) Close() {
	m.close()
}
