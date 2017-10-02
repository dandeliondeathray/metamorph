package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/dandeliondeathray/metamorph"
	"github.com/gorilla/mux"
)

func main() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	log.Println("Starting Metamorph...")
	server := metamorph.NewServer()

	kafka := metamorph.NewKafkaSystem(server)
	server.Kafka = kafka

	chSignal := make(chan os.Signal, 10)
	signal.Notify(chSignal, os.Interrupt)
	go stopSystemOnSignal(chSignal, kafka)

	server.Start()

	r := mux.NewRouter()
	r.HandleFunc("/", server.EventWebSocketHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":23572", r))
}

func stopSystemOnSignal(chSignal <-chan os.Signal, kafka *metamorph.KafkaSystem) {
	signal := <-chSignal
	log.Println("Handling signal:", signal)
	kafka.StopSystem()
	os.Exit(0)
}
