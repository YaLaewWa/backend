package server

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func (s *Server) initRoutes() {
	s.initSocket()
}

func (s *Server) initSocket() {
	s.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		s.db.Clients[c] = true
		defer func() { //Porperly close connection
			delete(s.db.Clients, c)
			c.Close()
		}()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv %v: %s", mt, msg)
			for k := range s.db.Clients {
				err = k.WriteMessage(mt, msg)
				if err != nil {
					break
				}
			}
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

}
