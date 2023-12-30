package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/harmony-cloud-live/server/internal/midi"
	"github.com/harmony-cloud-live/server/internal/osc"
	"github.com/harmony-cloud-live/server/internal/server"
)

func main() {
	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	midiPlayer, err := midi.NewMidiPlayer("IAC Driver Bus 1")
	if err != nil {
		panic(err)
	}
	defer midiPlayer.Close()
	
	oscClient := osc.NewOscClient("localhost", 8765)
	h := server.NewHarmonyCloudServer(midiPlayer, oscClient)

	s := &http.Server{
		Handler: h,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}