package metamorph

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type resetCompleteEvent struct {
	Type string `json:"type"`
}

type KafkaSystem struct {
	server    *Server
	zookeeper *kafkaSystemProcess
	kafka     *kafkaSystemProcess
	consumer  *consumer
	producer  *producer
}

func NewKafkaSystem(server *Server) *KafkaSystem {
	return &KafkaSystem{
		server:    server,
		zookeeper: newKafkaSystemProcess("Zookeeper", "/tmp/zookeeper", "/kafka/bin/zookeeper-server-start.sh", "/kafka/config/zookeeper.properties"),
		kafka:     newKafkaSystemProcess("Kafka", "/tmp/kafka-logs", "/kafka/bin/kafka-server-start.sh", "/kafka/config/server.properties"),
		consumer:  nil}
}

func (k *KafkaSystem) Reset(topics []string) {
	err := k.StopSystem()
	if err != nil {
		return
	}

	k.zookeeper.start()

	zookeeperProbe := newProbe("localhost", 2181, 2*time.Second)
	if !zookeeperProbe.run(10) {
		k.server.SendErrorEvent("Zookeeper process did not respond to probe within the timeout")
		return
	}

	k.kafka.start()

	kafkaProbe := newProbe("localhost", 9092, 2*time.Second)
	if !kafkaProbe.run(10) {
		k.server.SendErrorEvent("Kafka process did not respond to probe within the timeout")
		return
	}

	k.consumer = newConsumer(k.server)
	if err := k.consumer.start(topics); err != nil {
		log.Println("Could not start consumer:", err)
		k.server.SendErrorEvent("Could not start Kafka consumer")
		return
	}

	k.producer, err = newProducer()
	if err != nil {
		log.Println("Could not start producer:", err)
		k.server.SendErrorEvent("Could not start Kafka producer")
		return
	}

	k.sendResetCompleteEvent()
}

func (k *KafkaSystem) StopSystem() error {
	err := k.stopProcess(k.kafka)
	if err != nil {
		return err
	}

	err = k.stopProcess(k.zookeeper)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaSystem) Send(value, topic string) {
	err := k.producer.sendMessage(value, topic)
	if err != nil {
		errorMsg := fmt.Sprintf("Could not send message '%s' to topic '%s': %v", value, topic, err)
		log.Println(errorMsg)
		k.server.SendErrorEvent(errorMsg)
	}
}

func (k *KafkaSystem) stopProcess(process *kafkaSystemProcess) error {
	chExitSignal := make(chan bool)

	if process.isRunning() {
		err := process.stop(chExitSignal)
		if err != nil {
			log.Printf("Could not stop %s", process.name)
			k.server.SendErrorEvent("%s could not be stopped")
			return fmt.Errorf("Could not stop %s", process.name)
		}

		select {
		case <-chExitSignal:
			log.Printf("%s process exited...\n", process.name)
		case <-time.After(60 * time.Second):
			log.Printf("Wait for %s exit timed out...", process.name)
			k.server.SendErrorEvent("%s did not signal exit within timeout")
			return fmt.Errorf("Wait for %s exit timed out", process.name)
		}
	} else {
		log.Printf("%s is not running. Will not attempt to stop it.", process.name)
	}

	return nil
}

func (k *KafkaSystem) sendResetCompleteEvent() {
	k.server.SendEvent(resetCompleteEvent{Type: "reset_complete"})
}

func newKafkaSystemProcess(name string, pathToRemove string, cmd string, args ...string) *kafkaSystemProcess {
	return &kafkaSystemProcess{process: nil, name: name, cmd: cmd, args: args, pathToRemove: pathToRemove}
}

type kafkaSystemProcess struct {
	process      *os.Process
	name         string
	cmd          string
	args         []string
	pathToRemove string
}

func (z *kafkaSystemProcess) isRunning() bool {
	return z.process != nil
}

func (z *kafkaSystemProcess) start() {
	log.Printf("Starting %s...", z.name)

	os.RemoveAll(z.pathToRemove)

	cmd := exec.Command(z.cmd, z.args...)
	cmd.Start()
	z.process = cmd.Process
}

func (z *kafkaSystemProcess) stop(chExitSignal chan bool) error {
	err := z.process.Kill()
	if err != nil {
		return err
	}
	go z.signalWhenExited(chExitSignal)
	return nil
}

func (z *kafkaSystemProcess) signalWhenExited(chExitSignal chan bool) {
	_, err := z.process.Wait()
	if err != nil {
		log.Printf("Error when waiting for %s to exit: %v\n", z.name, err)
	}
	z.process = nil
	chExitSignal <- true
}
