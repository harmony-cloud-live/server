package music

import (
	"fmt"
	"strings"
)

func (a *ApiClient) BuildChord(songTitle string, chordSymbol string) (*Chord, error) {
	if songTitle == "" {
		return nil, fmt.Errorf("empty song name")
	}
	cache, ok := a.cache[strings.ReplaceAll(songTitle, " ", "")]
	if !ok {
		return nil, fmt.Errorf("song not found in cache: %s", songTitle)
	}

	chord, err := cache.getCleanChord(chordSymbol)
	if err != nil {
		return nil, err
	}

	chordSymbolInC, err := a.transposer.toC(chordSymbol, cache.key)
	if err != nil {
		return nil, err
	}
	chord.ChordSymbolInC = chordSymbolInC

	return chord, nil
} 

func (c *FallbackCache) getCleanChord(chordSymbol string) (*Chord, error) {
	pointers, ok := c.chordPointers[chordSymbol]
	if !ok {
		return nil, fmt.Errorf("chord not found: %s", chordSymbol)
	}

	chord := c.progressions[pointers[0].ProgressionIndex][pointers[0].ChordIndex]
	copy(c.chordPointers[chordSymbol], c.chordPointers[chordSymbol][1:])
	c.chordPointers[chordSymbol] = append(c.chordPointers[chordSymbol], pointers[0])
    
	cleanedChord, err := cleanChord(&chord)
	return cleanedChord, err
}
