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
	case GetNoteDelay:
		return h.marshalAndSend(ctx, c, GetNoteDelay, ControlPayload{NoteDelay: h.midiPlayer.GetNoteDelay()})
	case GetVelocity:
		return h.marshalAndSend(ctx, c, GetVelocity, ControlPayload{Velocity: h.midiPlayer.GetVelocity()})
	case GetPlaybackMode:
		return h.marshalAndSend(ctx, c, GetPlaybackMode, ControlPayload{PlaybackMode: h.state.GetPlaybackMode()})
	case GetManualModeChord:
		row, col := h.state.GetManualModeChord()
		return h.marshalAndSend(ctx, c, GetManualModeChord, ControlPayload{Row: row, Col: col})
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
			return fmt.Errorf("[ERROR] invalid username")
		}
	case NewIndex:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set index")
		}
		index, err := h.state.SetIndex(evt.Payload.Index)
		if err != nil {
			h.logf("[ERROR] setting index: %v", err)
			index, _ = h.state.SetIndex(0)
		}
		return h.marshalAndBroadcast(ctx, userId, NewIndex, ControlPayload{Index: index})
	case NewBeat:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set beat")
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
			return fmt.Errorf("[ERROR] only leader can set time signature")
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
				return fmt.Errorf("[ERROR] invalid user")
			}
			return h.setLeader(ctx, h.clients[newLeaderId])
		} else {
			return fmt.Errorf("[ERROR] invalid leaderId")
		}
	case NewLoop:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set loop")
		}
		start, end, err := h.state.SetLoop(evt.Payload.LoopStart, evt.Payload.LoopEnd)
		if err != nil {
			h.state.ClearLoop()
		} 
		return h.marshalAndBroadcast(ctx, userId, NewLoop, ControlPayload{LoopStart: start, LoopEnd: end})
	case SetNoteDelay:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set note delay")
		}
		noteDelay := h.midiPlayer.SetNoteDelay(evt.Payload.NoteDelay)
		return h.marshalAndBroadcast(ctx, userId, GetNoteDelay, ControlPayload{NoteDelay: noteDelay})
	case SetVelocity:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set velocity")
		}
		velocity := h.midiPlayer.SetVelocity(evt.Payload.Velocity)
		return h.marshalAndBroadcast(ctx, userId, GetVelocity, ControlPayload{Velocity: velocity})
	case SetPlaybackMode:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set playback mode")
		}
		playbackMode, err := h.state.SetPlaybackMode(evt.Payload.PlaybackMode)
		if err != nil {
			return err
		}
		return h.marshalAndBroadcast(ctx, userId, GetPlaybackMode, ControlPayload{PlaybackMode: playbackMode})
	case SetManualModeChord:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can set manual mode index")
		}
		row, col := h.state.SetManualModeChord(evt.Payload.Row, evt.Payload.Col)
		return h.marshalAndBroadcast(ctx, userId, GetManualModeChord, ControlPayload{Row: row, Col: col})
	case ManualChordDown:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can manually play chord")
		}
		if len(evt.Payload.Chords) != 1 {
			return fmt.Errorf("[ERROR] manual chord down requires exactly one chord")
		}
		songTitle := evt.Payload.SongTitle
		if songTitle == "" {
			songTitle = h.state.GetSongTitle()
		}
		chord, err := h.apiClient.BuildChord(songTitle, evt.Payload.Chords[0].ChordSymbol)
		if err != nil {
			return fmt.Errorf("[ERROR] building chord: %v", err)
		}
		h.midiPlayer.PlayChord(chord.MidiValues)
		h.oscClient.SendChord(*chord)
	case ManualChordUp:
		if h.leader != h.clients[userId] {
			return fmt.Errorf("[ERROR] only leader can manually change chord")
		}
		h.midiPlayer.StopAll()
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
		h.logf("[ERROR] caching playback state: %v", err)
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
	if h.leader != nil && h.leader.UserId == userId && username == "$listener$" {
		h.logf("[INFO] leader cannot become listener")
		return nil
	}
	h.clients[userId] = &Client{
		Conn: c,
		UserId: userId,
		Username: username,
	}
	if username == "$listener$" {
		h.logf("[INFO] listener connected")
		return nil
	}
	if h.leader == nil && h.cachedLeaderId != "" && h.clients[h.cachedLeaderId] != nil {
		h.logf("[INFO] auto set leader from cache: %v", h.clients[h.cachedLeaderId].Username)
		return h.setLeader(ctx, h.clients[h.cachedLeaderId])
	}
	if h.leader != nil && h.leader.UserId == userId && h.clients[userId] != nil {
		h.logf("[INFO] reconnecting leader: %v", h.clients[userId].Username)
		return h.setLeader(ctx, h.clients[userId])
	}
	return nil
}

func (h *HarmonyCloudServer) getClients() []*Client {
	var clients []*Client
	for _, c := range h.clients {
		if c.Username != "$listener$" {
			clients = append(clients, c)
		}
	}
	return clients
}
