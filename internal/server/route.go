package server

import (
	"log"
	"socket/internal/hub"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

func (s *Server) initRoutes() {
	s.initSocket()
	s.initAuth()
}

func (s *Server) initSocket() {
	s.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		id := uuid.New() //TODO: change to user id in the future
		hubChannel := make(chan []byte)
		payload := &hub.RegisterPayload{
			Channel: hubChannel,
			ID:      id,
		}
		s.messageHub.Register <- payload //register new connection to hub

		defer func() { //Porperly close connection
			s.messageHub.Unregister <- id
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
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			s.db.Message = append(s.db.Message, msg)
			s.messageHub.Broadcast <- msg

			msg = <-hubChannel
			err = c.WriteMessage(1, msg)
			if err != nil {
				break
			}
		}
	}))

}

func (s *Server) initAuth() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.handler.userHandler.Register)
	authRoutes.Post("/login", s.handler.userHandler.Login)
}
