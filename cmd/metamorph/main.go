package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dandeliondeathray/metamorph"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Metamorph...")
	server := metamorph.NewServer()
	server.Start()

	r := mux.NewRouter()
	r.HandleFunc("/", server.EventWebSocketHandler)

	http.Handle("/", r)

	go sendFakeEvents(server)
	log.Fatal(http.ListenAndServe(":23572", r))
}

type fakeEvent struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Count int    `json:"count"`
}

func sendFakeEvents(s *metamorph.Server) {
	count := 0
	for {
		s.SendEvent(fakeEvent{"message", "gurka", count})
		count++
		time.Sleep(500 * time.Millisecond)
	}
}
