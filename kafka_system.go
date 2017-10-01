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
		zookeeper: newKafkaSystemProcess("Zookeeper", "/bin/sleep", "10000"),
		kafka:     newKafkaSystemProcess("Kafka", "/bin/sleep", "100000")}
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

	k.zookeeper.start()

	// TODO: Wait for Zookeeper to start properly

	k.kafka.start()

	// TODO: Await Kafka up and running

	k.sendResetCompleteEvent()
}

func (k *KafkaSystem) stopProcess(process *kafkaSystemProcess) error {
	chExitSignal := make(chan bool)

	if k.zookeeper.isRunning() {
		err := process.stop(chExitSignal)
		if err != nil {
			log.Printf("Could not stop %s", process.name)
			k.server.SendErrorEvent("%s could not be stopped")
			return fmt.Errorf("Could not stop Zookeeper")
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
