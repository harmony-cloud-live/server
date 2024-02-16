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
	"github.com/harmony-cloud-live/server/internal/music"
	"github.com/harmony-cloud-live/server/internal/osc"
	"github.com/harmony-cloud-live/server/internal/server"
	"github.com/redis/go-redis/v9"
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
		return err
	}
	defer midiPlayer.Close()
	
	oscClient := osc.NewOscClient("192.168.1.70", 7000)

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	transposer, err := music.NewTransposer("internal/data/transpositions.csv")
	if err != nil {
		return err
	}
	
	apiClient, err := music.NewApiClient("192.168.1.212:5000", "internal/data/chords.json", transposer)
	if err != nil {
		return err
	}

	h, err := server.NewHarmonyCloudServer(midiPlayer, oscClient, rdb, apiClient)
	if err != nil {
		return err
	}

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
