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
	fmt.Println(unescaped)
	lines := bufio.NewScanner(strings.NewReader(unescaped))
	
	for lines.Scan() {
		if !foundChords {
			foundChords = strings.Contains(lines.Text(), "CHORDS") 
		} else {
			line := strings.Split(strings.TrimSpace(lines.Text()), ":")
			if len(line) != 2 {
				if strings.Contains(lines.Text(), "]") {
					break
				} else {
					continue
				}
			}
			chordSymbol := line[0]
			notes := strings.Split(strings.TrimSpace(line[1]), " ")

			var midiValues []uint8
			for _, note := range notes {
				val, err := strconv.Atoi(note)
				if err != nil {
					if len(chords) < 10 {
						return nil, fmt.Errorf("failed to convert note to int: %v", err)
					}
				}
				midiValues = append(midiValues, uint8(val))
			}

			chords = append(chords, Chord{ChordSymbol: chordSymbol, MidiValues: midiValues})
		}
	}
	
	return chords, nil
}