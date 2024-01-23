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
		h.logf("get clients %s", h.clients[userId].Username)
		return h.marshalAndSend(ctx, c, GetClients, h.getClients())
	case GetLeader:
		return h.marshalAndSend(ctx, c, GetLeader, h.getLeader())
	case SetUsername:
		if username, ok := evt.Payload.(string); ok {
			if h.clients[userId] == nil || h.clients[userId].Username != username || h.clients[userId].Conn != c {
				h.logf("set username: %v", evt.Payload)
				err := h.addClient(ctx, userId, username, c)
				if err != nil {
					return err
				}
				return h.marshalAndBroadcast(ctx, userId, GetClients, h.getClients())
			}
		} else {
			return fmt.Errorf("invalid username")
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
		} else {
			return fmt.Errorf("invalid index")
		}
	case NewBeat:
		h.logf("new beat: %v", evt.Payload)
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set beat")
		}
		if newBeat, ok := evt.Payload.(int); ok {
			h.currentBeat = newBeat
			return h.marshalAndBroadcast(ctx, userId, NewBeat, h.currentBeat)
		} else {
			return fmt.Errorf("invalid beat")
		}
	case NewMainSequence:
		h.logf("new main sequence: %v", evt.Payload)
		return h.newMainSequence(ctx)
	case NewLeader:
		h.logf("new leader: %v", evt.Payload)
		if newLeaderId, ok := evt.Payload.(string); ok {
			if h.clients[newLeaderId] == nil {
				return fmt.Errorf("invalid user")
			}
			return h.setLeader(ctx, h.clients[newLeaderId])
		} else {
			return fmt.Errorf("invalid leaderId")
		}
	}
	return err
}

func (h *HarmonyCloudServer) newMainSequence(ctx context.Context) error {
    h.currentIndex = 0
    h.currentBeat = 0
    h.mainSequence = getDummyData()

	err := h.marshalAndBroadcast(ctx, "", NewMainSequence, h.mainSequence)
	if err != nil {
		return err
	}
	h.cacheMainSequence(ctx)
	return nil
}

func (h *HarmonyCloudServer) setLeader(ctx context.Context, newLeader *Client) error {
	h.mu.Lock()
	h.leader = newLeader
	h.mu.Unlock()

	err := h.marshalAndBroadcast(ctx, "", GetLeader, h.getLeader())
	if err != nil {
		return err
	}
	h.cacheLeaderId(ctx)
	return nil
}

func (h *HarmonyCloudServer) getLeader() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.leader == nil {
		return ""
	}
	return fmt.Sprintf("%v|%v", h.leader.UserId, h.leader.Username)
}

func (h *HarmonyCloudServer) addClient(ctx context.Context, userId string, username string, c *websocket.Conn) error {
	h.clients[userId] = &Client{
		Conn: c,
		UserId: userId,
		Username: username,
	}
	if h.leader == nil && h.cachedLeaderId != "" && h.clients[h.cachedLeaderId] != nil {
		h.logf("auto set leader from cache: %v", h.clients[h.cachedLeaderId])
		return h.setLeader(ctx, h.clients[h.cachedLeaderId])
	}
	if h.leader != nil && h.leader.UserId == userId {
		h.logf("reconnecting leader: %v", h.clients[userId])
		return h.setLeader(ctx, h.clients[userId])
	}
	return nil
}

func (h *HarmonyCloudServer) getClients() []*ClientData {
	var clients []*ClientData
	for _, c := range h.clients {
		clients = append(clients, &ClientData{UserId: c.UserId, Username: c.Username})
	}
	return clients
}
