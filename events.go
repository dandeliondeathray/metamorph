package metamorph

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

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
		_, _, err := c.ReadMessage()
		if err != nil {
			return
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
