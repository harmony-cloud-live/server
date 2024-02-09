package server

import (
	"context"
	"fmt"

	"nhooyr.io/websocket"
)

func (h *HarmonyCloudServer) handleMidiEvent(ctx context.Context, c *websocket.Conn, userId string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, rawMsg, err := c.Read(ctx)
	if err != nil {
		return err
	}

	evt, err := unmarshalMidiEvent(rawMsg)
	if err != nil {
		return err
	}
	
	switch evt.Type {
	case ChordDown:
		if evt.Index < 0 || evt.Index >= len(*h.mainSequence) {
			return fmt.Errorf("invalid index")
		}
		chord := (*h.mainSequence)[evt.Index]
		h.midiPlayer.PlayChord(chord.MidiValues)
		h.oscClient.SendNotes(chord.MidiValues)
		if err != nil {
			return err
		}
		h.oscClient.SendChordSymbol(chord.ChordSymbolInC)
	case ChordUp:
		if evt.Index < 0 || evt.Index >= len(*h.mainSequence) {
			return fmt.Errorf("invalid index")
		}
		h.midiPlayer.StopChord((*h.mainSequence)[evt.Index].MidiValues)
		h.oscClient.SendRelease()
	case StopAll:
		h.midiPlayer.StopAll()
		h.oscClient.SendRelease()
	}

	return err
}
