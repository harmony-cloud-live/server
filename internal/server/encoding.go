package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func marshalControlEvent(eventType ControlEventType, payload ControlPayload) ([]byte, error) {
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    size := 1 + 2 + len(payloadBytes)
    buffer := bytes.NewBuffer(make([]byte, 0, size))

    buffer.WriteByte(byte(eventType))

    if err := binary.Write(buffer, binary.LittleEndian, uint16(len(payloadBytes))); err != nil {
        return nil, err
    }

    buffer.Write(payloadBytes)

    return buffer.Bytes(), nil
}

func unmarshalControlEvent(buffer []byte) (*ControlEvent, error) {
    if len(buffer) < 3 {
        return nil, fmt.Errorf("buffer too short")
    }

    eventType := ControlEventType(buffer[0])

    payloadLength := binary.LittleEndian.Uint16(buffer[1:3])

    if int(payloadLength) > len(buffer)-3 {
        return nil, fmt.Errorf("buffer is shorter than expected")
    }

    var payload ControlPayload
    if err := json.Unmarshal(buffer[3:3+payloadLength], &payload); err != nil {
        return nil, fmt.Errorf("failed to unmarshal payload: %v", err)
    }

	if num, ok := payload.(float64); ok {
		return &ControlEvent{Type: eventType, Payload: int(num)}, nil
	}

    return &ControlEvent{Type: eventType, Payload: payload}, nil
}

func unmarshalMidiEvent(data []byte) (*MidiEvent, error) {
    if len(data) < 1 {
        return nil, fmt.Errorf("data is too short")
    }

    eventType := MidiEventType(data[0])

    if eventType == StopAll || len(data) < 2 {
        return &MidiEvent{Type: eventType, Index: 0}, nil
    }

    return &MidiEvent{Type: eventType, Index: int(data[1])}, nil
}