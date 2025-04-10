package server

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
)

func (s *Server) initRoutes() {
	s.initSocket()
}

func (s *Server) initSocket() {

	s.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		fmt.Println(c.Locals("Host"))
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

}
