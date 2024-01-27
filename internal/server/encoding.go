package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"nhooyr.io/websocket"
)

func marshalControlEvent(eventType ControlEventType, payload ControlPayload) ([]byte, error) {
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("marshal error: %v", err)
    }

    size := 1 + 2 + len(payloadBytes)
    buffer := bytes.NewBuffer(make([]byte, 0, size))

    buffer.WriteByte(byte(eventType))

    if err := binary.Write(buffer, binary.LittleEndian, uint16(len(payloadBytes))); err != nil {
        return nil, fmt.Errorf("marshal error: %v", err)
    }

    buffer.Write(payloadBytes)

    return buffer.Bytes(), nil
}

func unmarshalControlEvent(buffer []byte) (*ControlEvent, error) {
    if len(buffer) < 3 {
        return nil, fmt.Errorf("unmarshal error: buffer too short")
    }

    eventType := ControlEventType(buffer[0])

    payloadLength := binary.LittleEndian.Uint16(buffer[1:3])

    if int(payloadLength) > len(buffer)-3 {
        return nil, fmt.Errorf("unmarshal error: buffer is shorter than expected")
    }

    var payload ControlPayload
    if err := json.Unmarshal(buffer[3:3+payloadLength], &payload); err != nil {
        return nil, fmt.Errorf("unmarshal error: %v", err)
    }

    return &ControlEvent{Type: eventType, Payload: payload}, nil
}

func unmarshalMidiEvent(buffer []byte) (*MidiEvent, error) {
    if len(buffer) < 1 {
        return nil, fmt.Errorf("unmarshal error: buffer is too short")
    }

    eventType := MidiEventType(buffer[0])

    if eventType == StopAll || len(buffer) < 2 {
        return &MidiEvent{Type: eventType, Index: 0}, nil
    }

    return &MidiEvent{Type: eventType, Index: int(buffer[1])}, nil
}

func (h *HarmonyCloudServer) marshalAndSend(ctx context.Context, 
	c *websocket.Conn, eventType ControlEventType, payload ControlPayload) error {
	msg, err := marshalControlEvent(eventType, payload)
	if err != nil {
		return err
	}
	return c.Write(ctx, websocket.MessageBinary, msg)
}

func (h *HarmonyCloudServer) marshalAndBroadcast(ctx context.Context, 
	userId string, eventType ControlEventType, payload ControlPayload) error {
	msg, err := marshalControlEvent(eventType, payload)
	if err != nil {
		return err
	}
	return h.broadcast(ctx, userId, msg)
}