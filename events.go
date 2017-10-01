package metamorph

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type generalCommand struct {
	Type string `json:"type"`
}

type Server struct {
	connections map[*websocket.Conn]bool
	chEvents    chan interface{}
	mutex       sync.Mutex
}

func NewServer() *Server {
	return &Server{make(map[*websocket.Conn]bool), make(chan interface{}, 100), sync.Mutex{}}
}

func (s *Server) SendEvent(ev interface{}) {
	s.chEvents <- ev
}

func (s *Server) Start() {
	go s.eventsWriter()
}

func (s *Server) EventWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	s.newConnection(c)
	go s.eventsReader(c)
}

func (s *Server) eventsWriter() {
	for event := range s.chEvents {
		log.Printf("Broadcasting event %v\n", event)

		s.mutex.Lock()

		for c := range s.connections {
			c.WriteJSON(event)
		}

		s.mutex.Unlock()
	}
}

func (s *Server) eventsReader(c *websocket.Conn) {
	defer c.Close()
	defer s.closeConnection(c)

	for {
		_, p, err := c.ReadMessage()
		log.Printf("Read message from Metamorph client: %v", p)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var genCommand generalCommand
		err = json.Unmarshal(p, &genCommand)
		if err != nil {
			log.Println("Parse command type error:", err)
			break
		}
		log.Printf("Command: %v", genCommand)

		if genCommand.Type == "reset_kafka_system" {
			// TODO: Perform reset
		}
	}
}

func (s *Server) newConnection(c *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Println("New connection established.")

	s.connections[c] = true
}

func (s *Server) closeConnection(c *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Println("Closing connection...")

	delete(s.connections, c)
}
