package server

import (
	"context"

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
		chord, err := h.state.GetChord(evt.Index)
		if err != nil {
			return err
		}
		go h.oscClient.SendChord(chord)
		h.midiPlayer.PlayChord(chord.MidiValues)
	case ChordUp:
		// chord, err := h.state.GetChord(evt.Index)
		// if err != nil {
		// 	return err
		// }
		// h.midiPlayer.StopChord(chord.MidiValues)
		// h.oscClient.SendRelease()
		fallthrough
	case StopAll:
		h.midiPlayer.StopAll()
		h.oscClient.SendRelease()
	}

	return err
}
