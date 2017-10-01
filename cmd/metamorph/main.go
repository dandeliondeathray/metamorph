package main

import (
	"log"
	"net/http"
	"os"

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

	server.Start()

	r := mux.NewRouter()
	r.HandleFunc("/", server.EventWebSocketHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":23572", r))
}
