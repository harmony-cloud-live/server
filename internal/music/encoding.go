package music

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type RequestBody struct {
	Input   string `json:"input"`
	Length  int    `json:"length"`
}

type ResponseBody struct {
    Message string `json:"message"`
    Status  string `json:"status"`
}

func marshalRequest(collection ChordCollection, startingChord Chord, length int) (*bytes.Reader, error) {
	input := "[KEY: " + string(collection.Key) + "\n COLLECTION:"
	for _, chord := range collection.ChordSymbols {
		input += " " + chord
	}
	input += "\\n CHORDS:\\n "
	if startingChord.ChordSymbol != "" {
		input += startingChord.ChordSymbol + ":"
		for _, midiValue := range startingChord.MidiValues {
			input += " " + strconv.Itoa(int(midiValue))
		}
	}
	body := RequestBody{
		Input:  input,
		Length: length,
	}

	enc, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(enc), nil
}

func unmarshalResponse(resp *http.Response) ([]Chord, error) {
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var body ResponseBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
    }
	
	if body.Status != "success" {
		return nil, fmt.Errorf("failed to generate music: %s", body.Message)
	}

	chords := make([]Chord, 0)
	foundChords := false
	
	unescaped := strings.ReplaceAll(body.Message, "\\n", "\n")
	lines := bufio.NewScanner(strings.NewReader(unescaped))
	
	for lines.Scan() {
		if !foundChords {
			foundChords = strings.Contains(lines.Text(), "CHORDS") 
		} else {
			line := strings.TrimSpace(lines.Text())
			if strings.Contains(line, "]") {
				break
			}

			chord, notes, found := strings.Cut(line, ": ")
			if !found {
				break
			}

			var midiValues []uint8
			for _, note := range strings.Split(notes, " ") {
				val, err := strconv.Atoi(note)
				if err != nil {
					return nil, fmt.Errorf("failed to convert note to int: %v", err)
				}
				midiValues = append(midiValues, uint8(val))
			}

			chords = append(chords, Chord{ChordSymbol: chord, MidiValues: midiValues})
		}
	}
	
	if len(chords) < 10 {
		return nil, fmt.Errorf("failed to generate music: %s", body.Message)
	}
	
	return chords[:len(chords)-1], nil
}

func CleanChords(chords []Chord) ([]Chord, error) {
	var result []Chord
	for _, chord := range chords {
		if !chord.isValid() {
			return nil, fmt.Errorf("invalid chord: %v", chord)
		}
		cleanedChord, err := cleanChord(&chord)
		if err != nil {
			return nil, err
		}
		result = append(result, *cleanedChord)
	}
	return result, nil
}	

func cleanChord(chord *Chord) (*Chord, error) {
	seen := make(map[uint8]bool)
	cleanedChord := Chord{ChordSymbol: chord.ChordSymbol}
	for _, note := range chord.MidiValues {
		if _, ok := seen[note]; !ok {
			seen[note] = true
			cleanedChord.MidiValues = append(cleanedChord.MidiValues, note)
		}
	}
	return &cleanedChord, nil
}