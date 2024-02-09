package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type ApiClient struct {
	baseUrl 		string
	client 			*http.Client
	transposer 		*Transposer
	collectionMap   map[string]ChordCollection
	cache 			map[string]*SongCache
}

type ChordCollection struct {
	Title 			string `json:"title"`
	ChordSymbols 	[]string `json:"chordSymbols"`
	Key 			KeySignature `json:"key"`
}

func NewApiClient(baseUrl string, chordFile string, transposer *Transposer) (*ApiClient, error) {
	client := &http.Client{
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
		if collection.Title == "" {
			return nil, fmt.Errorf("invalid title: %s", collection.Title)
		}
		if !isValidKeySignature(collection.Key) {
			return nil, fmt.Errorf("invalid key: %s", collection.Key)
		}
		if len(collection.ChordSymbols) == 0 {
			return nil, fmt.Errorf("invalid chord collection: %s", collection.ChordSymbols)
		}
		collectionMap[collection.Title] = collection
	}
	
	cache, err := getDummyData()
	if err != nil {
		fmt.Printf("[WARNING] failed to get dummy data: %v", err)
	}

	return &ApiClient{
		baseUrl: baseUrl,
		client: client,
		transposer: transposer,
		collectionMap: collectionMap,
		cache: cache,
	}, nil
}

func (a *ApiClient) GenerateMusicWithFallback(songName string, startingChord Chord, length int) *[]Chord {
	chords, err := a.GenerateMusic(songName, startingChord, length)
	if err != nil {
		fmt.Printf("[WARNING] failed to generate music, falling back to cache: %v \n", err)
		return a.fallbackToCache(songName, startingChord)
	}
	return &chords
}

func (a *ApiClient) GenerateMusic(songName string, startingChord Chord, length int) ([]Chord, error) {
	collection, ok := a.collectionMap[songName]
	if !ok {
		return nil, fmt.Errorf("invalid song name: %s", songName)
	}

	body, err := marshalRequest(collection, startingChord, length)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Post("http://" + a.baseUrl + "/generate", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	chords, err := unmarshalResponse(resp)
	if err != nil {
		return nil, err
	}

	err = a.transposer.PopulateTranspositions(&chords, collection.Key)
	if err != nil {
		return nil, err
	}
	return chords, nil
}

func (a *ApiClient) fallbackToCache(songName string, startingChord Chord) *[]Chord {
	var result []Chord

	normalizedName := strings.ReplaceAll(songName, " ", "")
	if songName != "" && a.cache[normalizedName] != nil {
		song := a.cache[normalizedName]
		if startingChord.ChordSymbol != "" {
			result = song.GetProgressionByChord(startingChord.ChordSymbol)
		} else {
			result = song.GetProgression()
		}
	} else {
		for _, song := range a.cache {
			result = song.GetProgression()
			break
		}
	}
	return &result
}

