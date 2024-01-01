package server

import (
	"context"
	"fmt"

	"nhooyr.io/websocket"
)

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
	case GetIndex:
		return h.marshalAndSend(ctx, c, NewIndex, h.currentIndex)
	case GetBeat:
		return h.marshalAndSend(ctx, c, NewBeat, h.currentBeat)
	case GetMainSequence:
		return h.marshalAndSend(ctx, c, NewMainSequence, h.mainSequence)
	case GetClients:
		h.logf("get clients", h.clients[userId])
		return h.marshalAndSend(ctx, c, GetClients, h.getClients())
	case GetLeader:
		return h.marshalAndSend(ctx, c, GetLeader, h.getLeader())
	case SetUsername:
		if username, ok := evt.Payload.(string); ok {
			if h.clients[userId] == nil || h.clients[userId].Username != username {
				h.logf("set username: %v", evt.Payload)
				h.addClient(userId, username, c)
				return h.marshalAndBroadcast(ctx, userId, GetClients, h.getClients())
			}
		}
	case NewIndex:
		h.logf("new index: %v", evt.Payload)
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set index")
		}
		if newIndex, ok := evt.Payload.(int); ok {
			if newIndex < 0 || newIndex >= len(h.mainSequence) {
				return fmt.Errorf("invalid index")
			}
			h.currentIndex = newIndex
			err := h.marshalAndBroadcast(ctx, userId, NewIndex, h.currentIndex)
			if err != nil {
				return err
			}
			chord := h.mainSequence[h.currentIndex]
			h.oscClient.SendNotes(chord.MidiValues)
			h.oscClient.SendChordSymbol(chord.ChordSymbol)
		}
	case NewBeat:
		h.logf("new beat: %v", evt.Payload)
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set beat")
		}
		if newBeat, ok := evt.Payload.(int); ok {
			h.currentBeat = newBeat
			return h.marshalAndBroadcast(ctx, userId, NewBeat, h.currentBeat)
		}
	case NewMainSequence:
		h.logf("new main sequence: %v", evt.Payload)
		h.mainSequence = getDummyData()
		return h.marshalAndBroadcast(ctx, userId, NewMainSequence, h.mainSequence)
	case NewLeader:
		h.logf("new leader: %v", evt.Payload)
		if newLeaderId, ok := evt.Payload.(string); ok {
			if h.clients[newLeaderId] == nil {
				return fmt.Errorf("invalid user")
			}
			return h.setLeader(ctx, h.clients[newLeaderId])
		}
	}

	return err
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

func (h *HarmonyCloudServer) setLeader(ctx context.Context, newLeader *Client) error {
	h.leader = newLeader
	// do we need to get senderId
	return h.marshalAndBroadcast(ctx, "", GetLeader, h.getLeader())
}

func (h *HarmonyCloudServer) getLeader() string {
	if h.leader == nil {
		return ""
	}
	return fmt.Sprintf("%v|%v", h.leader.UserId, h.leader.Username)
}
