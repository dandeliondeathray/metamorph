package metamorph

import (
	"fmt"
	"log"
	"net"
	"time"
)

type probe struct {
	port  int
	host  string
	delay time.Duration
}

func newProbe(host string, port int, delay time.Duration) *probe {
	return &probe{port: port, host: host, delay: delay}
}

func (p *probe) run(count int) bool {
	for i := 0; i < count; i++ {
		time.Sleep(p.delay)
		if p.doProbing() {
			return true
		}
	}
	return false
}

func (p *probe) doProbing() bool {
	address := fmt.Sprintf("%s:%d", p.host, p.port)
	log.Println("Probing address", address)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err == nil {
		conn.Close()
		log.Println("Probe successful.")
		return true
	}
	log.Println("Probe not successful.")
	return false
}
