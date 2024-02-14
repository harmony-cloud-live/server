package music

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"

	"github.com/agnivade/levenshtein"
)


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

type SymbolsInC map[string]string

type Transposer struct {
	keyToSymbolsInC map[KeySignature]SymbolsInC 
}

func NewTransposer(source string) (*Transposer, error) {
	t := Transposer{
		keyToSymbolsInC: make(map[KeySignature]SymbolsInC),
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
	
	keys := make([]KeySignature, len(records[0]))
	for i, tonic := range records[0] {
		key := KeySignature(tonic)
		if !key.isValid() {
			return nil, fmt.Errorf("invalid tonic found in CSV: %s", tonic)
		}
		t.keyToSymbolsInC[key] = make(SymbolsInC)
		keys[i] = key
	}
	
	for _, record := range records {
		for i, symbol := range record {
			t.keyToSymbolsInC[keys[i]][symbol] = record[0]
		}
	}
	
	return &t, nil
}

func (t *Transposer) PopulateTranspositions(progression *[]Chord, key KeySignature) error {
	if !key.isValid() {
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
	fmt.Printf("[INFO] populated %d transpositions\n", len(*progression))
	return nil
}

func (t *Transposer) toC(chord string, key KeySignature) (string, error) {
	symbolsInC, ok := t.keyToSymbolsInC[key]
	if !ok {
		return "", fmt.Errorf("invalid key: %s", key)
	}

	transposedChord, ok := symbolsInC[chord]
	if !ok {
		fmt.Println("Chord", chord, "not found in key", key, "attempting to find nearest chord")
		nearest := t.nearestSymbol(chord, key)
		transposedChord, ok = symbolsInC[nearest]
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
	for symbol := range t.keyToSymbolsInC[key] {
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