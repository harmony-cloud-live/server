package osc

import (
	// "fmt"
	"strconv"
	"strings"

	"github.com/harmony-cloud-live/server/internal/music"
	"github.com/hypebeast/go-osc/osc"
)

type OscClient struct {
	client *osc.Client
}

func NewOscClient(ip string, port int) *OscClient {
	client := osc.NewClient(ip, port)
	return &OscClient{
		client: client,
	}
}

func (o *OscClient) SendChord(chord music.Chord) {
	o.sendChordSymbol(chord.ChordSymbolInC)
	o.sendNotes(chord.MidiValues)
}

func (o *OscClient) sendNotes(notes []uint8) {
	var builder strings.Builder
	for i, note := range notes {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(strconv.Itoa(int(note)))
	}

	msg := osc.NewMessage("/notes", builder.String())
	o.client.Send(msg)
}

func (o *OscClient) sendChordSymbol(chord string) {
	msg := osc.NewMessage("/chord", chord)
	o.client.Send(msg)
}

func (o *OscClient) SendRelease() {
	msg := osc.NewMessage("/release")
	o.client.Send(msg)
}
