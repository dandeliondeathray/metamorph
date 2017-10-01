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
}

func NewKafkaSystem(server *Server) *KafkaSystem {
	return &KafkaSystem{
		server:    server,
		zookeeper: newKafkaSystemProcess("Zookeeper", "/kafka/bin/zookeeper-server-start.sh", "/kafka/config/zookeeper.properties"),
		kafka:     newKafkaSystemProcess("Kafka", "/kafka/bin/kafka-server-start.sh", "/kafka/config/server.properties")}
}

func (k *KafkaSystem) Reset() {
	err := k.stopProcess(k.kafka)
	if err != nil {
		return
	}

	err = k.stopProcess(k.zookeeper)
	if err != nil {
		return
	}

	time.Sleep(30 * time.Second)

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

	k.sendResetCompleteEvent()
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

func newKafkaSystemProcess(name string, cmd string, args ...string) *kafkaSystemProcess {
	return &kafkaSystemProcess{process: nil, name: name, cmd: cmd, args: args}
}

type kafkaSystemProcess struct {
	process *os.Process
	name    string
	cmd     string
	args    []string
}

func (z *kafkaSystemProcess) isRunning() bool {
	return z.process != nil
}

func (z *kafkaSystemProcess) start() {
	log.Printf("Starting %s...", z.name)
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
