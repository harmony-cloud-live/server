package data

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Chord struct {
    ChordSymbol string `json:"chordSymbol"`
    ChordSymbolInC string
    MidiValues  []uint8
}

type ChordPointer struct {
	ProgressionIndex int
	ChordIndex       int
}

type SongCache struct {
	key KeySignature
	progressions [][]Chord
	chordPointers map[string][]ChordPointer
}

func (c *SongCache) GetKey() KeySignature {
	return c.key
}

func (c *SongCache) GetProgression() []Chord {
	totalProgressions := len(c.progressions)
	return c.progressions[rand.Intn(totalProgressions)][:24]
}

func (c *SongCache) GetProgressionByChord(chordSymbol string) []Chord {
	pointers := c.chordPointers[chordSymbol]
	if len(pointers) == 0 {
		return c.GetProgression()
	}
	copy(c.chordPointers[chordSymbol], c.chordPointers[chordSymbol][1:])
	c.chordPointers[chordSymbol] = append(c.chordPointers[chordSymbol], pointers[0])
	progression := c.progressions[pointers[0].ProgressionIndex][pointers[0].ChordIndex:]

	if len(progression) < 24 {
		return progression
	} else {
		return progression[:24]
	}
}

func getDummyData() (map[string]*SongCache, error) {
	dirs, err := os.ReadDir("internal/data")
	if err != nil {
		return nil, fmt.Errorf("failed to list data dirs: %s", err)
	}
	
	songCaches := make(map[string]*SongCache)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dirNameParts := strings.Split(dir.Name(), "_")
		if len(dirNameParts) != 2 {
			fmt.Println("Unexpected song directory format:", dir.Name())
			continue
		}
		songName := dirNameParts[0]
		key := dirNameParts[1]
		if !isValidKeySignature(KeySignature(key)) {
			fmt.Println("Invalid key signature:", key)
			continue
		}
		cache, err := processSongDirectory("internal/data/" + dir.Name())
		if err != nil {
			fmt.Println("Error processing song directory:", dir, err)
			continue
		}
		cache.key = KeySignature(key)
		songCaches[songName] = cache
		fmt.Println("Loaded fallback cache:", dir.Name())
	}
	
	return songCaches, nil
}

func processSongDirectory(path string) (*SongCache, error) {
	c := SongCache{
		progressions: make([][]Chord, 0),
		chordPointers: make(map[string][]ChordPointer),
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %s", err)
	}
	
	for i, file := range files {
		if filepath.Ext(file.Name()) != ".txt" {
			continue
		}
		
		chords, err := parseFile(path + "/" + file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to parse file: %s", err)
		}
		
		c.progressions = append(c.progressions, chords)

		end := len(chords)
		if len(chords) > 24 {
			end = len(chords) - 24
		}

		for j, chord := range chords[:end] {
			c.chordPointers[chord.ChordSymbol] = append(c.chordPointers[chord.ChordSymbol], ChordPointer{
				ProgressionIndex: i, 
				ChordIndex: j,
			})
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

func parseFile(filename string) ([]Chord, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}
	
	lines := strings.Split(string(content), "\n")

	var chords []Chord
	for _, line := range lines {
		if line == "" {
			continue
		}
		chordParts := strings.SplitN(line, ": ", 2)
		if len(chordParts) != 2 {
			fmt.Println("Unexpected line format:", line)
			continue
		}
		chordSymbol := chordParts[0]
		midiValuesStr := strings.Split(chordParts[1], " ")

		var midiValues []uint8
		for _, valStr := range midiValuesStr {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Println("Error converting MIDI value:", err)
				continue
			}
			midiValues = append(midiValues, uint8(val))
		}

		chords = append(chords, Chord{
			ChordSymbol: chordSymbol,
			MidiValues:  midiValues,
		})
	}

	return chords, nil
}
