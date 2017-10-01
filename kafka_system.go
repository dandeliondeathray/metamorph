package metamorph

type resetCompleteEvent struct {
	Type string `json:"type"`
}

type KafkaSystem struct {
	server *Server
}

func NewKafkaSystem(server *Server) *KafkaSystem {
	return &KafkaSystem{server: server}
}

func (k *KafkaSystem) Reset() {
	// TODO: Stop currently running Kafka system, if any

	// TODO: Start new Kafka system

	// TODO: Await Kafka up and running

	// TODO: Send reset complete event
}

func (k *KafkaSystem) SendResetCompleteEvent() {
	k.server.SendEvent(resetCompleteEvent{Type: "reset_complete"})
}
