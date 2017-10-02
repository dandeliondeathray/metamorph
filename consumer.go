package metamorph

import "github.com/Shopify/sarama"
import "log"

type consumer struct {
	server             *Server
	consumer           sarama.Consumer
	partitionConsumers []sarama.PartitionConsumer
}

func newConsumer(server *Server) *consumer {
	return &consumer{server, nil, make([]sarama.PartitionConsumer, 5)}
}

func (c *consumer) start(topics []string) error {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		return err
	}

	for i := range topics {
		err = c.consumeTopic(consumer, topics[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *consumer) consumeTopic(consumer sarama.Consumer, topic string) error {
	log.Println("Creating consumer for topic", topic)
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for i := range partitions {
		p := partitions[i]
		log.Printf("Creating partition consumer for topic %s and partition %d", topic, p)
		partitionConsumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetOldest)
		if err != nil {
			log.Printf("ERROR: Couldn't create partition consumer %d", p)
			return err
		}
		go c.forwardMessage(topic, partitionConsumer.Messages())
		c.partitionConsumers = append(c.partitionConsumers, partitionConsumer)
	}

	return nil
}

type Message struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Value []byte `json:"message"`
}

func (c *consumer) forwardMessage(topic string, chMessages <-chan *sarama.ConsumerMessage) {
	for m := range chMessages {
		log.Println("Received message:", m)
		c.server.SendEvent(Message{"message", topic, m.Value})
	}
}
