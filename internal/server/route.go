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

		for index := range s.db.Message { // send old message to the new user
			err := c.WriteMessage(1, s.db.Message[index])
			if err != nil {
				log.Println("read:", err)
				break
			}
		}

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			s.db.Message = append(s.db.Message, msg)
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
