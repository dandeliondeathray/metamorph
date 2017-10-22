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
	metamorphAddress := os.Getenv("METAMORPH_ADDRESS")
	if len(metamorphAddress) == 0 {
		metamorphAddress = ":23572"
	}

	kafkaRoot := os.Getenv("METAMORPH_KAFKA_ROOT")
	log.Println("Using Kafka root", kafkaRoot)

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	log.Println("Starting Metamorph...")
	server := metamorph.NewServer()

	kafkaConfig := metamorph.DefaultKafkaConfig()
	kafkaConfig.Root = kafkaRoot
	kafka := metamorph.NewKafkaSystem(server, kafkaConfig)
	server.Kafka = kafka

	chSignal := make(chan os.Signal, 10)
	signal.Notify(chSignal, os.Interrupt)
	go stopSystemOnSignal(chSignal, kafka)

	server.Start()

	r := mux.NewRouter()
	r.HandleFunc("/", server.EventWebSocketHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(metamorphAddress, r))
}

func stopSystemOnSignal(chSignal <-chan os.Signal, kafka *metamorph.KafkaSystem) {
	signal := <-chSignal
	log.Println("Handling signal:", signal)
	kafka.StopSystem()
	os.Exit(0)
}
