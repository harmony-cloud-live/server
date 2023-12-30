package osc

import "github.com/hypebeast/go-osc/osc"

type OscClient struct {
	client *osc.Client
}

func NewOscClient(ip string, port int) *OscClient {
	client := osc.NewClient(ip, port)
	return &OscClient{
		client: client,
	}
}

func (o *OscClient) SendNotes(notes []uint8) {
	for _, note := range notes {
		msg := osc.NewMessage("/notes")
		msg.Append(int32(note))
		o.client.Send(msg)
	}
}

func (o *OscClient) SendChordSymbol(chord string) {
	msg := osc.NewMessage("/chord", chord)
	o.client.Send(msg)
}

func (o *OscClient) SendRelease() {
	msg := osc.NewMessage("/release")
	o.client.Send(msg)
}