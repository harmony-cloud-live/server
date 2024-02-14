package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harmony-cloud-live/server/internal/music"
)


func (h *HarmonyCloudServer) cacheLeaderId(ctx context.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.leader != nil {
        err := h.rdb.Set(ctx, "leaderId", h.leader.UserId, 0).Err()
        if err != nil {
            h.logf("[ERROR] caching leaderId: %v", err)
        }
	}
}

func (h *HarmonyCloudServer) getCachedLeaderId(ctx context.Context) string {
	leaderId, err := h.rdb.Get(ctx, "leaderId").Result()
	if err != nil {
        h.logf("[ERROR] getting cached leaderId: %v", err)
		return ""
	}
	return leaderId
}

func (h *HarmonyCloudServer) cachePlaybackState(ctx context.Context, state *music.PlaybackState) error {
    chords := state.GetChords()
    keySignature := state.GetKeySignature()
    songTitle := state.GetSongTitle()

    var chordSymbols []string
    var chordSymbolsInC []string
    var midiValues [][]uint8
    for _, chord := range chords {
        chordSymbols = append(chordSymbols, chord.ChordSymbol)
        chordSymbolsInC = append(chordSymbolsInC, chord.ChordSymbolInC)
        midiValues = append(midiValues, chord.MidiValues)
    }

    encChordSymbols, err := json.Marshal(chordSymbols)
    if err != nil {
        return fmt.Errorf("error marshalling chordSymbols: %v", err)
    }

    encChordSymbolsInC, err := json.Marshal(chordSymbolsInC)
    if err != nil {
        return fmt.Errorf("error marshalling chordSymbolsInC: %v", err)
    }

    encMidiValues, err := json.Marshal(midiValues)
    if err != nil {
        return fmt.Errorf("error marshalling midiValues: %v", err)
    }
    
    pipe := h.rdb.TxPipeline()
    pipe.Set(ctx, "mainSequenceSongTitle", songTitle, 0)
    pipe.Set(ctx, "mainSequenceKeySignature", keySignature, 0)
    pipe.Set(ctx, "mainSequenceChordSymbols", encChordSymbols, 0)
    pipe.Set(ctx, "mainSequenceChordSymbolsInC", encChordSymbolsInC, 0)
    pipe.Set(ctx, "mainSequenceMidiValues", encMidiValues, 0)

    _, err = pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("error caching playbackState: %v", err)
    }
    return nil
}

func (h *HarmonyCloudServer) getCachedPlaybackState(ctx context.Context) (*music.PlaybackState, error) {
    var chordSymbols []string
    var chordSymbolsInC []string
    var midiValues [][]uint8

    pipe := h.rdb.TxPipeline()
    cachedSongTitle := pipe.Get(ctx, "mainSequenceSongTitle")
    cachedKeySignature := pipe.Get(ctx, "mainSequenceKeySignature")
    cachedChordSymbols := pipe.Get(ctx, "mainSequenceChordSymbols")
    cachedChordSymbolsInC := pipe.Get(ctx, "mainSequenceChordSymbolsInC")
    cachedMidiValues := pipe.Get(ctx, "mainSequenceMidiValues")

    _, err := pipe.Exec(ctx)
    if err != nil {
        return nil, err
    }

    songTitle := cachedSongTitle.Val()
    keySignature := music.KeySignature(cachedKeySignature.Val())

    err = json.Unmarshal([]byte(cachedChordSymbols.Val()), &chordSymbols)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling chordSymbols: %v", err)
    }

    err = json.Unmarshal([]byte(cachedChordSymbolsInC.Val()), &chordSymbolsInC)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling chordSymbolsInC: %v", err)
    }

    err = json.Unmarshal([]byte(cachedMidiValues.Val()), &midiValues)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling midiValues: %v", err)
    }

    var chords []music.Chord
    for i, symbol := range chordSymbols {
        chords = append(chords, music.Chord{
            ChordSymbol: symbol, 
            ChordSymbolInC: chordSymbolsInC[i], 
            MidiValues: midiValues[i],
        })
    }
    
    state := music.NewPlaybackState()
    err = state.SetMainSequence(&music.ChordProgression{
        Title: songTitle,
        Chords: chords,
        Key: keySignature,
    })

    if err != nil {
        return nil, err
    }

    h.logf("[INFO] loaded playbackState from cache: %s", state.GetSongTitle())
    return state, nil
}
