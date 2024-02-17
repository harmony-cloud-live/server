package music

import "fmt"

type PlaybackState struct {
	mainSequence *ChordProgression
	timeSignature TimeSignature
	index int
	beat int
	loopStart int
	loopEnd int
    playbackMode string
    manualModeRow int
    manualModeCol int
}

func NewPlaybackState() *PlaybackState {
	return &PlaybackState{
		timeSignature: TimeSignature{Upper: 4, Lower: 4},
		index: 0,
		beat: 0,
		loopStart: -1,
		loopEnd: -1,
	}
}

func (p *PlaybackState) GetManualModeChord() (int, int) {
    return p.manualModeRow, p.manualModeCol
}

func (p *PlaybackState) SetManualModeChord(row, col int) (int, int) {
    p.manualModeRow = row
    p.manualModeCol = col
    return p.manualModeRow, p.manualModeCol
}

func (p *PlaybackState) GetPlaybackMode() string {
    return p.playbackMode
}

func (p *PlaybackState) SetPlaybackMode(playbackMode string) (string, error) {
    if playbackMode != "ai" && playbackMode != "manual" {
        return "", fmt.Errorf("invalid playback mode")
    }
    p.playbackMode = playbackMode
    return p.playbackMode, nil
}

func (p *PlaybackState) GetSongTitle() string {
    if p.mainSequence == nil {
        return ""
    }
    return p.mainSequence.Title
}

func (p *PlaybackState) GetChords() []Chord {
    if (p.mainSequence == nil) {
        return []Chord{}
    }
    return p.mainSequence.Chords
}

func (p *PlaybackState) SetMainSequence(mainSequence *ChordProgression) error {
    if mainSequence == nil || !mainSequence.isValid() {
        return fmt.Errorf("invalid mainSequence")
    }
    p.mainSequence = mainSequence
    return nil
}

func (p *PlaybackState) IsValidMainSequence() bool {
    return p.mainSequence != nil && p.mainSequence.isValid()
}

func (p *PlaybackState) GetChord(index int) (Chord, error) {
    if !p.IsValidMainSequence() {
        return Chord{}, fmt.Errorf("invalid main sequence")
    }
    if !p.IsValidIndex(index) {
        return Chord{}, fmt.Errorf("invalid index: %d", index)
    }
    return p.mainSequence.Chords[index], nil
}

func (p *PlaybackState) GetIndex() int {
    return p.index
}

func (p *PlaybackState) SetIndex(index int) (int, error) {
    if p.IsValidIndex(index) {
        p.index = index
        return p.index, nil
    }
    return 0, fmt.Errorf("invalid index: %d", index)
}

func (p *PlaybackState) IsValidIndex(index int) bool {
    if (p.mainSequence == nil) {
        return false
    }
    return index >= 0 && index < len(p.mainSequence.Chords)
}

func (p *PlaybackState) GetBeat() int {
    return p.index
}

func (p *PlaybackState) SetBeat(beat int) (int, error) {
    if p.IsValidBeat(beat) {
        p.beat = beat
        return p.beat, nil
    }
    return 0, fmt.Errorf("invalid beat: %d", beat)
}

func (p *PlaybackState) IsValidBeat(beat int) bool {
    return beat >= 0 && beat < p.timeSignature.Upper
}

func (p *PlaybackState) SetLoop(start int, end int) (int, int, error) {
    if !p.IsValidIndex(start) || !p.IsValidIndex(end) || start >= end {
        return -1, -1, fmt.Errorf("invalid loop (%v, %v)", start, end)
    }
    p.loopStart = start
    p.loopEnd = end
    return p.loopStart, p.loopEnd, nil
}

func (p *PlaybackState) GetLoop() (int, int) {
    return p.loopStart, p.loopEnd
}

func (p *PlaybackState) ClearLoop() {
    p.loopStart = -1
    p.loopEnd = -1
}

func (p *PlaybackState) GetTimeSignature() TimeSignature {
    return p.timeSignature
}

func (p *PlaybackState) SetTimeSignature(timeSignature TimeSignature) (TimeSignature, error) {
    if timeSignature.isValid() {
        p.timeSignature = timeSignature
        return p.timeSignature, nil
    }
    return p.timeSignature, fmt.Errorf("invalid time signature: %v", timeSignature)
}

func (p *PlaybackState) GetKeySignature() string {
    return string(p.mainSequence.Key)
}
