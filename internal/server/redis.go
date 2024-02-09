package server

import (
	"context"
	"encoding/json"

	"github.com/harmony-cloud-live/server/internal/data"
)


func (h *HarmonyCloudServer) cacheLeaderId(ctx context.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.leader != nil {
        err := h.rdb.Set(ctx, "leaderId", h.leader.UserId, 0).Err()
        if err != nil {
            h.logf("error caching leaderId: %v", err)
        }
	}
}

func (h *HarmonyCloudServer) getCachedLeaderId(ctx context.Context) string {
	leaderId, err := h.rdb.Get(ctx, "leaderId").Result()
	if err != nil {
		return ""
	}
	return leaderId
}

func (h *HarmonyCloudServer) cacheMainSequence(ctx context.Context) {
    var chordSymbols []string
    var midiValues [][]uint8
    for _, chord := range *h.mainSequence {
        chordSymbols = append(chordSymbols, chord.ChordSymbol)
        midiValues = append(midiValues, chord.MidiValues)
    }

    encChordSymbols, err := json.Marshal(chordSymbols)
    if err != nil {
        h.logf("error marshalling chordSymbols: %v", err)
        return
    }

    encMidiValues, err := json.Marshal(midiValues)
    if err != nil {
        h.logf("error marshalling midiValues: %v", err)
        return
    }
    
    pipe := h.rdb.TxPipeline()
    pipe.Set(ctx, "mainSequenceSongName", h.songName, 0)
    pipe.Set(ctx, "mainSequenceKey", string(h.key), 0)
    pipe.Set(ctx, "mainSequenceChordSymbols", encChordSymbols, 0)
    pipe.Set(ctx, "mainSequenceMidiValues", encMidiValues, 0)

    _, err = pipe.Exec(ctx)
    if err != nil {
        h.logf("error caching mainSequence: %v", err)
    }
}

func (h *HarmonyCloudServer) getCachedMainSequence(ctx context.Context) *[]data.Chord {
    var mainSequence []data.Chord
    var chordSymbols []string
    var midiValues [][]uint8

    pipe := h.rdb.TxPipeline()
    encChordSymbols := pipe.Get(ctx, "mainSequenceChordSymbols")
    encMidiValues := pipe.Get(ctx, "mainSequenceMidiValues")
    songName := pipe.Get(ctx, "mainSequenceSongName")
    key := pipe.Get(ctx, "mainSequenceKey")

    _, err := pipe.Exec(ctx)
    if err != nil {
        return nil
    }

    err = json.Unmarshal([]byte(encChordSymbols.Val()), &chordSymbols)
    if err != nil {
        h.logf("error unmarshalling chordSymbols: %v", err)
        return nil
    }

    err = json.Unmarshal([]byte(encMidiValues.Val()), &midiValues)
    if err != nil {
        h.logf("error unmarshalling midiValues: %v", err)
        return nil
    }

    for i, symbol := range chordSymbols {
        mainSequence  = append(mainSequence, data.Chord{ChordSymbol: symbol, MidiValues: midiValues[i]})
    }

    if songName.Val() == "" || key.Val() == "" {
        return nil
    }
    
    h.songName = songName.Val()
    h.key = data.KeySignature(key.Val())
    
    h.logf("loaded mainSequence from cache: %v in %v", h.songName, h.key)

    return &mainSequence
}
