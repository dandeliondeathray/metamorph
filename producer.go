package metamorph

import (
	"log"

	"github.com/Shopify/sarama"
)

type producer struct {
	syncProducer sarama.SyncProducer
}

func newProducer() (*producer, error) {
	p, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return nil, err
	}

	return &producer{syncProducer: p}, nil
}

func (p *producer) sendMessage(value, topic string) error {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Sent message: Partition %d, Offset %d", partition, offset)
	return nil
}
