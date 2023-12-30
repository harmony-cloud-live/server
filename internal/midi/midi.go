package midi

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiPlayer struct {
	send func(msg midi.Message) error
	close func()
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
}

func (m *MidiPlayer) StopAll() {
	for note := uint8(0); note < 128; note++ {
		m.send(midi.NoteOff(0, note))
	}
}

func (m *MidiPlayer) Close() {
	m.close()
}