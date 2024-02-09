package data

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"

	"github.com/agnivade/levenshtein"
)

type KeySignature string

const (
	C  KeySignature = "C"
	Db KeySignature = "Db"
	D  KeySignature = "D"
	Eb KeySignature = "Eb"
	E  KeySignature = "E"
	F  KeySignature = "F"
	Gb KeySignature = "Gb"
	G  KeySignature = "G"
	Ab KeySignature = "Ab"
	A  KeySignature = "A"
	Bb KeySignature = "Bb"
	B  KeySignature = "B"
)

func isValidKeySignature(key KeySignature) bool {
	switch key {
	case C, Db, D, Eb, E, F, Gb, G, Ab, A, Bb, B:
		return true
	default:
		return false
	}
}

var sharps = map[string]string{
	"C#": "Db",
	"D#": "Eb",
	"E#": "F",
	"F#": "Gb",
	"G#": "Ab",
	"A#": "Bb",
	"B#": "C",
}

var flats = map[string]string{
	"Cb": "B",
	"Db": "C#",
	"Eb": "D#",
	"Fb": "E",
	"Gb": "F#",
	"Ab": "G#",
	"Bb": "A#",
}

var naturals = map[string]string{
	"C": "B#",
	"B": "Cb",
	"E": "Fb",
	"F": "E#",
}

func getEnharmonic(chord string) string {
	if len(chord) > 1 && chord[1] == '#' {
		enharmonic, ok := sharps[chord[:2]] 
		if ok {
			return enharmonic + chord[2:]
		}
	} else if len(chord) > 1 && chord[1] == 'b' {
		enharmonic, ok := flats[chord[:2]] 
		if ok {
			return enharmonic + chord[2:]
		}
	} else if len(chord) > 0 {
		enharmonic, ok := naturals[string(chord[0])]
		if ok {
			return enharmonic + chord[1:]
		}
	} 
	return ""
}

type SymbolMap map[string]string

type Transposer struct {
	keyToSymbolMap map[KeySignature]SymbolMap
}

func NewTransposer(source string) (*Transposer, error) {
	t := Transposer{
		keyToSymbolMap: make(map[KeySignature]SymbolMap),
	}

	file, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read transpositions: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read transpositions: %v", err)
	}
	
	indexToKey := make(map[int]KeySignature)
	for i, key := range records[0] {
		ks := KeySignature(key)
		if !isValidKeySignature(ks) {
			return nil, fmt.Errorf("invalid key found in CSV: %s", key)
		}
		t.keyToSymbolMap[ks] = make(SymbolMap)
		indexToKey[i] = ks
	}
	
	for _, record := range records {
		for i, symbol := range record {
			t.keyToSymbolMap[indexToKey[i]][symbol] = record[0]
		}
	}
	
	return &t, nil
}

func (t *Transposer) PopulateTranspositions(progression *[]Chord, key KeySignature) error {
	if !isValidKeySignature(key) {
		return fmt.Errorf("invalid key: %s", key)
	}

	if progression == nil {
		return fmt.Errorf("nil progression")
	}

	for i, chord := range *progression {
		transposed, err := t.toC(chord.ChordSymbol, key)
		if err != nil {
			return fmt.Errorf("failed to transpose chord: %v", err)
		}
		(*progression)[i].ChordSymbolInC = transposed
	}
	fmt.Printf("Populated %d transpositions\n", len(*progression))
	return nil
}

func (t *Transposer) toC(chord string, key KeySignature) (string, error) {
	symbolMap, ok := t.keyToSymbolMap[key]
	if !ok {
		return "", fmt.Errorf("invalid key: %s", key)
	}

	transposedChord, ok := symbolMap[chord]
	if !ok {
		fmt.Println("Chord", chord, "not found in key", key, "attempting to find nearest chord")
		nearest := t.nearestSymbol(chord, key)
		transposedChord, ok = symbolMap[nearest]
		if !ok {
			return "", fmt.Errorf("nearest chord %s not found in key %s", chord, key)
		}
	}

	return transposedChord, nil
}

func (t *Transposer) nearestSymbol(chord string, key KeySignature) string {
	enharmonic := getEnharmonic(chord)
	minDistance := math.MaxInt
	minChord := ""
	for symbol := range t.keyToSymbolMap[key] {
		distance := levenshtein.ComputeDistance(symbol, chord)
		if distance < minDistance {
			minDistance = distance
			minChord = symbol
		}

		distance = levenshtein.ComputeDistance(symbol, enharmonic)
		if distance < minDistance {
			minDistance = distance
			minChord = symbol
		}
	}
	fmt.Println("Chord:", chord, "Nearest chord:", minChord, "distance:", minDistance)
	return minChord
}