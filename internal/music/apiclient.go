package music

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ApiClient struct {
	logf 			func(f string, v ...interface{})
	baseUrl 		string
	httpClient 		*http.Client
	transposer 		*Transposer
	collectionMap   map[string]ChordCollection
	cache 			map[string]*FallbackCache
}

func NewApiClient(baseUrl string, chordFile string, transposer *Transposer) (*ApiClient, error) {
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}

	file, err := os.Open(chordFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read chord file: %v", err)
	}
	defer file.Close()

	var collections []ChordCollection
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&collections)
	if err != nil {
		return nil, fmt.Errorf("failed to decode chord file: %v", err)
	}
	
	var collectionMap = make(map[string]ChordCollection)
	for _, collection := range collections {
		if collection.isValid() {
			collectionMap[collection.Title] = collection
		} else {
			return nil, fmt.Errorf("invalid chord collection: %v", collection)
		}
	}
	
	cache, err := loadFallbackData()
	if err != nil {
		return nil, fmt.Errorf("failed to load fallback data: %v", err)
	}

	return &ApiClient{
		logf: log.Printf,
		baseUrl: baseUrl,
		httpClient: httpClient,
		transposer: transposer,
		collectionMap: collectionMap,
		cache: cache,
	}, nil
}

func (a *ApiClient) GenerateMusicWithFallback(songTitle string, startingChord Chord, length int) (*ChordProgression, error) {
	chordProgression, err := a.GenerateMusic(songTitle, startingChord, length)
	if err != nil {
		a.logf("[WARNING] failed to generate music, falling back to cache: %v", err)
		chordProgression, err = a.fallbackToCache(songTitle, startingChord)
		if err != nil {
			return nil, fmt.Errorf("fallback cache produced invalid chord progression: %v", err)
		}
	}

	return chordProgression, nil
}

func (a *ApiClient) GenerateMusic(songTitle string, startingChord Chord, length int) (*ChordProgression, error) {
	collection, ok := a.collectionMap[songTitle]
	if !ok {
		return nil, fmt.Errorf("invalid song title: %s", songTitle)
	}

	body, err := marshalRequest(collection, startingChord, length)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := a.httpClient.Post("http://" + a.baseUrl + "/generate", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	a.logf("[INFO] elapsed: %.2f secs", time.Since(start).Seconds())

	chords, err := unmarshalResponse(resp)
	if err != nil {
		return nil, err
	}

	err = a.transposer.PopulateTranspositions(&chords, collection.Key)
	if err != nil {
		return nil, err
	}

	chords, err = cleanChords(chords)
	if err != nil {
		return nil, err
	}

	return &ChordProgression{
		Title:  songTitle,
		Chords: chords,
		Key:    collection.Key,
	}, nil
}

func (a *ApiClient) fallbackToCache(songTitle string, startingChord Chord) (*ChordProgression, error) {
	var chords []Chord
	var key KeySignature

	if songTitle != "" {
		song, ok := a.cache[strings.ReplaceAll(songTitle, " ", "")]
		if !ok {
			return nil, fmt.Errorf("song not found in cache: %s", songTitle)
		}
		key = song.key

		if startingChord.ChordSymbol != "" {
			chords = song.GetProgressionByChord(startingChord.ChordSymbol)
		} else {
			chords = song.GetProgression()
		}
	} else {
		return nil, fmt.Errorf("empty song name")
	}

	err := a.transposer.PopulateTranspositions(&chords, key)
	if err != nil {
		return nil, err
	}
	
	chords, err = cleanChords(chords)
	if err != nil {
		return nil, err
	}

	return &ChordProgression{
		Title:  songTitle,
		Chords: chords,
		Key:    key,
	}, nil
}

