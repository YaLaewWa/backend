package server

import (
	"github.com/gofiber/contrib/websocket"
)

func (s *Server) initRoutes() {
	s.initSocket()
	s.initAuth()
	s.initMessage()
}

func (s *Server) initSocket() {
	s.app.Get("/ws", websocket.New(s.handler.socketMessageHandler.InitConnection))

}

func (s *Server) initAuth() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.handler.userHandler.Register)
	authRoutes.Post("/login", s.handler.userHandler.Login)
}

func (s *Server) initMessage() {
	messageRoutes := s.app.Group("/messages")
	messageRoutes.Get("/", s.handler.socketMessageHandler.GetAll)
}
