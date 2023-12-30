package server

import (
	"context"
	"fmt"

	"nhooyr.io/websocket"
)

func (h *HarmonyCloudServer) broadcast(ctx context.Context, msg []byte) error {
	for _, u := range h.clients {
		err := u.Conn.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HarmonyCloudServer) handleControlEvent(ctx context.Context, c *websocket.Conn, userId string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, rawMsg, err := c.Read(ctx)
	if err != nil {
		return err
	}
	
	evt, err := unmarshalControlEvent(rawMsg)
	if err != nil {
		h.logf("failed to unmarshal control event: %v", err)
		return err
	}

	switch evt.Type {
	case SetUsername:
		h.logf("set username: %v", evt.Payload)
		h.addClient(userId, evt.Payload.(string), c)
		msg, err := marshalControlEvent(GetClients, h.getClients())
		if err != nil {
			return err
		}
		h.broadcast(ctx, msg)
	case GetClients:
		msg, err := marshalControlEvent(GetClients, h.getClients())
		if err != nil {
			return err
		}
		err = c.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			return err
		}
	case NewLeader:
		h.logf("new leader: %v", evt.Payload)
		newLeader := h.clients[evt.Payload.(string)]
		if newLeader == nil {
			return fmt.Errorf("invalid leader")
		}
		h.leader = newLeader
		payload := fmt.Sprintf("%v|%v", h.leader.UserId, h.leader.Username)
		msg, err := marshalControlEvent(GetLeader, payload)
		if err != nil {
			return err
		}
		h.broadcast(ctx, msg)
	case GetLeader:
		if h.leader == nil {
			return fmt.Errorf("invalid leader")
		}
		payload := fmt.Sprintf("%v|%v", h.leader.UserId, h.leader.Username)
		msg, err := marshalControlEvent(GetLeader, payload)
		if err != nil {
			return err
		}
		err = c.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			return err
		}
	case NewIndex:
		h.logf("new index: %v", evt.Payload)
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set index")
		}
		h.currentIndex = evt.Payload.(int)
		msg, err := marshalControlEvent(NewIndex, h.currentIndex)
		if err != nil {
			return err
		}
		h.broadcast(ctx, msg)

		chord := h.mainSequence[h.currentIndex]
		h.oscClient.SendNotes(chord.MidiValues)
		h.oscClient.SendChordSymbol(chord.ChordSymbol)
	case GetIndex:
		msg, err := marshalControlEvent(NewIndex, h.currentIndex)
		if err != nil {
			return err
		}
		err = c.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			return err
		}
	case NewBeat:
		h.logf("new beat: %v", evt.Payload)
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set beat")
		}
		h.currentBeat = evt.Payload.(int)
		msg, err := marshalControlEvent(NewBeat, h.currentBeat)
		if err != nil {
			return err
		}
		h.broadcast(ctx, msg)
	case GetBeat:
		msg, err := marshalControlEvent(NewBeat, h.currentBeat)
		if err != nil {
			return err
		}
		err = c.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			return err
		}
	case NewMainSequence:
		h.logf("new main sequence: %v", evt.Payload)
		h.mainSequence = getDummyData()
		msg, err := marshalControlEvent(NewMainSequence, h.mainSequence)
		if err != nil {
			return err
		}
		h.broadcast(ctx, msg)
	case GetMainSequence:
		msg, err := marshalControlEvent(NewMainSequence, h.mainSequence)
		if err != nil {
			return err
		}
		err = c.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			return err
		}
	}

	return err
}


