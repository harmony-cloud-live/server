package server

import (
	"context"
	"fmt"

	"github.com/harmony-cloud-live/server/internal/music"
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
		return h.marshalAndSend(ctx, c, NewIndex, ControlPayload{Index: h.state.GetIndex()})
	case GetBeat:
		return h.marshalAndSend(ctx, c, NewBeat, ControlPayload{Beat: h.state.GetBeat()})
	case GetMainSequence:
		return h.marshalAndSend(ctx, c, NewMainSequence, ControlPayload{
			Chords: h.state.GetChords(),
			SongTitle: h.state.GetSongTitle(),
		})
	case GetClients:
		return h.marshalAndSend(ctx, c, GetClients, ControlPayload{Clients: h.getClients()})
	case GetLeader:
		return h.marshalAndSend(ctx, c, NewLeader, ControlPayload{LeaderId: h.getLeader()})
	case GetTimeSignature:
		return h.marshalAndSend(ctx, c, NewTimeSignature, ControlPayload{TimeSignature: h.state.GetTimeSignature()})
	case GetLoop:
		loopStart, loopEnd := h.state.GetLoop()
		return h.marshalAndSend(ctx, c, NewLoop, ControlPayload{LoopStart: loopStart, LoopEnd: loopEnd})
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
		index, err := h.state.SetIndex(evt.Payload.Index)
		if err != nil {
			return err
		}
		return h.marshalAndBroadcast(ctx, userId, NewIndex, ControlPayload{Index: index})
	case NewBeat:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set beat")
		}
		beat, err := h.state.SetBeat(evt.Payload.Beat)
		if err != nil {
			return err
		}
		return h.marshalAndBroadcast(ctx, userId, NewBeat, ControlPayload{Beat: beat})
	case NewMainSequence:
		return h.newMainSequence(ctx, evt.Payload.SongTitle)
	case NewTimeSignature:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("only leader can set time signature")
		}
		timeSignature, err := h.state.SetTimeSignature(evt.Payload.TimeSignature)
		if err != nil {
			return err
		} 
		return h.marshalAndBroadcast(ctx, "", NewTimeSignature, ControlPayload{TimeSignature: timeSignature})
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
		start, end, err := h.state.SetLoop(evt.Payload.LoopStart, evt.Payload.LoopEnd)
		if err != nil {
			h.state.ClearLoop()
		} 
		return h.marshalAndBroadcast(ctx, userId, NewLoop, ControlPayload{LoopStart: start, LoopEnd: end})
	}
	return err
}

func (h *HarmonyCloudServer) newMainSequence(ctx context.Context, songTitle string) error {
	var startingChord music.Chord
	index := h.state.GetIndex()
	if index != 0 && h.state.IsValidIndex(index) && h.state.GetSongTitle() == songTitle {
		currentChord, err := h.state.GetChord(index)
		if err != nil {
			return err
		}
		startingChord = currentChord
	}

	mainSequence, err := h.apiClient.GenerateMusicWithFallback(songTitle, startingChord, 400)
	if err != nil {
		return err
	}

	err = h.state.SetMainSequence(mainSequence)
	if err != nil {
		return err
	}

	err = h.marshalAndBroadcast(ctx, "", NewMainSequence, ControlPayload{
		Chords: h.state.GetChords(),
		SongTitle: h.state.GetSongTitle(),
	})
	if err != nil {
		return err
	}

	err = h.cachePlaybackState(ctx, h.state)
	if err != nil {
		h.logf("error caching playback state: %v", err)
	}
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

func (h *HarmonyCloudServer) getClients() []*Client {
	var clients []*Client
	for _, c := range h.clients {
		clients = append(clients, c)
	}
	return clients
}
