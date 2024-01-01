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
		h.logf("failed to unmarshal control event: %v", err)
		return err
	}
	
	switch evt.Type {
	case ChordDown:
		if evt.Index >= len(h.mainSequence) {
			return fmt.Errorf("invalid index")
		}

		h.midiPlayer.PlayChord(h.mainSequence[evt.Index].MidiValues)
	case StopAll:
		h.midiPlayer.StopAll()
		h.oscClient.SendRelease()
	}

	return err
}

