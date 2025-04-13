package server

import (
	_ "socket/docs"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/swagger"
)

func (s *Server) initRoutes() {
	s.initSocket()
	s.initAuth()
	s.initSwagger()
}

func (s *Server) initSocket() {
	s.app.Get("/ws", websocket.New(s.handler.socketMessageHandler.InitConnection))

}

func (s *Server) initAuth() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.handler.userHandler.Register)
	authRoutes.Post("/login", s.handler.userHandler.Login)
}

func (s *Server) initSwagger() {
	s.app.Get("/swagger/*", swagger.HandlerDefault) // default

	s.app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		DeepLinking:  false,
		DocExpansion: "none",
	}))

}
