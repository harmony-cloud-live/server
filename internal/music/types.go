package music

type Chord struct {
    ChordSymbol string `json:"chordSymbol"`
    ChordSymbolInC string
    MidiValues  []uint8
}

func (c Chord) isValid() bool {
	return c.ChordSymbol != "" && c.ChordSymbolInC != "" && len(c.MidiValues) != 0
}

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

func (k KeySignature) isValid() bool {
	switch k {
	case C, Db, D, Eb, E, F, Gb, G, Ab, A, Bb, B:
		return true
	default:
		return false
	}
}

type TimeSignature struct {
    Upper int `json:"upper"`
    Lower int `json:"lower"`
}

func (t TimeSignature) isValid() bool {
    return t.Upper > 0 && t.Lower > 0
}

type ChordProgression struct {
	Title 	string
	Chords 	[]Chord
	Key 	KeySignature
}

func (c ChordProgression) isValid() bool {
    return len(c.Chords) > 0 && c.Key.isValid() && c.Title != ""
}

type ChordCollection struct {
	Title 			string `json:"title"`
	ChordSymbols 	[]string `json:"chordSymbols"`
	Key 			KeySignature `json:"key"`
}

func (c ChordCollection) isValid() bool {
	return c.Title != "" && c.Key.isValid() && len(c.ChordSymbols) != 0
}