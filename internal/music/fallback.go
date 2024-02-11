package music

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ChordPointer struct {
	ProgressionIndex int
	ChordIndex       int
}

type FallbackCache struct {
	key KeySignature
	progressions [][]Chord
	chordPointers map[string][]ChordPointer
}

const defaultLength = 24

func (c *FallbackCache) GetProgression() []Chord {
	totalProgressions := len(c.progressions)
	return c.progressions[rand.Intn(totalProgressions)][:defaultLength]
}

func (c *FallbackCache) GetProgressionByChord(chordSymbol string) []Chord {
	pointers := c.chordPointers[chordSymbol]
	if len(pointers) == 0 {
		return c.GetProgression()
	}
	copy(c.chordPointers[chordSymbol], c.chordPointers[chordSymbol][1:])
	c.chordPointers[chordSymbol] = append(c.chordPointers[chordSymbol], pointers[0])
	progression := c.progressions[pointers[0].ProgressionIndex][pointers[0].ChordIndex:]

	if len(progression) < defaultLength {
		return progression
	} else {
		return progression[:defaultLength]
	}
}

func loadFallbackData() (map[string]*FallbackCache, error) {
	dirs, err := os.ReadDir("internal/data")
	if err != nil {
		return nil, fmt.Errorf("failed to list data dirs: %s", err)
	}
	
	caches := make(map[string]*FallbackCache)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		songTitle, key, found := strings.Cut(dir.Name(), "_")
		if !found {
			return nil, fmt.Errorf("unexpected fallback directory format: %s", dir.Name())
		}
		if !KeySignature(key).isValid() {
			return nil, fmt.Errorf("invalid key signature: %s", key)
		}
		cache, err := parseFallbackDirectory("internal/data/" + dir.Name())
		if err != nil {
			return nil, fmt.Errorf("error processing song directory: %s %s", dir.Name(), err)
		}
		cache.key = KeySignature(key)
		caches[songTitle] = cache
	}
	
	return caches, nil
}

func parseFallbackDirectory(path string) (*FallbackCache, error) {
	c := FallbackCache{
		progressions: make([][]Chord, 0),
		chordPointers: make(map[string][]ChordPointer),
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %s", err)
	}
	
	for i, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {
			chords, err := parseFallbackFile(path + "/" + file.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to parse file: %s", err)
			}
			c.progressions = append(c.progressions, chords)

			end := len(chords)
			if len(chords) > defaultLength {
				end = len(chords) - defaultLength
			}

			for j, chord := range chords[:end] {
				ptr := ChordPointer{
					ProgressionIndex: i, 
					ChordIndex: j,
				}
				c.chordPointers[chord.ChordSymbol] = append(c.chordPointers[chord.ChordSymbol], ptr)
			}
		}
	}
	
	for k, v := range c.chordPointers {
		rand.Shuffle(len(v), func(i, j int) {
			v[i], v[j] = v[j], v[i]
		})
		c.chordPointers[k] = v
	}

	return &c, nil
}

func parseFallbackFile(filename string) ([]Chord, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}
	lines := strings.Split(string(content), "\n")

	var chords []Chord
	for _, line := range lines {
		chord, midi, found := strings.Cut(line, ": ")
		if !found {
			continue
		}

		midiStrings := strings.Split(midi, " ")

		var midiValues []uint8
		for _, midiString := range midiStrings  {
			val, err := strconv.Atoi(midiString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse midi value: %s", midiString)
			}
			midiValues = append(midiValues, uint8(val))
		}

		chords = append(chords, Chord{
			ChordSymbol: chord,
			MidiValues:  midiValues,
		})
	}

	return chords, nil
}
