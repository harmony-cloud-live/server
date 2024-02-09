package server

import (
	"context"
	"fmt"

	"github.com/harmony-cloud-live/server/internal/data"
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
		return h.marshalAndSend(ctx, c, NewIndex, ControlPayload{Index: h.currentIndex})
	case GetBeat:
		return h.marshalAndSend(ctx, c, NewBeat, ControlPayload{Beat: h.currentBeat})
	case GetMainSequence:
		return h.marshalAndSend(ctx, c, NewMainSequence, ControlPayload{Chords: *h.mainSequence})
	case GetClients:
		return h.marshalAndSend(ctx, c, GetClients, ControlPayload{Clients: h.getClients()})
	case GetLeader:
		return h.marshalAndSend(ctx, c, NewLeader, ControlPayload{LeaderId: h.getLeader()})
	case GetTimeSignature:
		return h.marshalAndSend(ctx, c, NewTimeSignature, ControlPayload{TimeSignature: h.timeSignature})
	case GetLoop:
		return h.marshalAndSend(ctx, c, NewLoop, ControlPayload{LoopStart: h.loopStart, LoopEnd: h.loopEnd})
	case SetUsername:
		username := evt.Payload.Username
		if username != "" {
			if h.clients[userId] == nil || h.clients[userId].Username != username || h.clients[userId].Conn != c {
				err := h.addClient(ctx, userId, username, c)
				if err != nil {
					return err
				}
				return h.marshalAndBroadcast(ctx, userId, GetClients, ControlPayload{Clients: h.getClients()})
			}
		} else {
			return fmt.Errorf("invalid username")
		}
	case NewIndex:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set index")
		}
		newIndex := evt.Payload.Index
		if newIndex >= 0 && newIndex < len(*h.mainSequence) {
			h.currentIndex = newIndex
			err := h.marshalAndBroadcast(ctx, userId, NewIndex, ControlPayload{Index: h.currentIndex})
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid index")
		}
	case NewBeat:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set beat")
		}
		newBeat := evt.Payload.Beat
		if newBeat >= 0 && newBeat < int(h.timeSignature.Upper) {
			h.currentBeat = newBeat
			return h.marshalAndBroadcast(ctx, userId, NewBeat, ControlPayload{Beat: h.currentBeat})
		} else {
			return fmt.Errorf("invalid beat")
		}
	case NewMainSequence:
		fmt.Println("new main sequence", evt.Payload.SongName)
		return h.newMainSequence(ctx, evt.Payload.SongName)
	case NewTimeSignature:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set time signature")
		}
		newTimeSignature := evt.Payload.TimeSignature
		if newTimeSignature.isValid() {
			h.timeSignature = newTimeSignature
			return h.marshalAndBroadcast(ctx, "", NewTimeSignature, ControlPayload{TimeSignature: h.timeSignature})
		} else {
			return fmt.Errorf("invalid time signature")
		}
	case NewLeader:
		newLeaderId := evt.Payload.LeaderId
		if newLeaderId != "" {
			if h.clients[newLeaderId] == nil {
				return fmt.Errorf("invalid user")
			}
			return h.setLeader(ctx, h.clients[newLeaderId])
		} else {
			return fmt.Errorf("invalid leaderId")
		}
	case NewLoop:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set loop")
		}
		newLoopStart := evt.Payload.LoopStart
		newLoopEnd := evt.Payload.LoopEnd
		if newLoopStart >= 0 && newLoopEnd >= 0 && newLoopStart < newLoopEnd && newLoopEnd < len(*h.mainSequence) {
			h.loopStart = newLoopStart
			h.loopEnd = newLoopEnd
			return h.marshalAndBroadcast(ctx, userId, NewLoop, ControlPayload{LoopStart: h.loopStart, LoopEnd: h.loopEnd})
		} else {
			h.loopStart = -1
			h.loopEnd = -1
			return h.marshalAndBroadcast(ctx, userId, NewLoop, ControlPayload{LoopStart: h.loopStart, LoopEnd: h.loopEnd})
		}
	}
	return err
}

func (h *HarmonyCloudServer) newMainSequence(ctx context.Context, songName string) error {
	startingChord := data.Chord{}
	if h.currentIndex != 0 {
		startingChord = (*h.mainSequence)[h.currentIndex]
	} 

	h.mainSequence = h.apiClient.GenerateMusicWithFallback(songName, startingChord, 400)
	h.currentIndex = 0
    h.currentBeat = 0

	err := h.marshalAndBroadcast(ctx, "", NewMainSequence, ControlPayload{Chords: *h.mainSequence})
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

	err := h.marshalAndBroadcast(ctx, "", GetLeader, ControlPayload{LeaderId: h.leader.UserId})
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
	return h.leader.UserId
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
