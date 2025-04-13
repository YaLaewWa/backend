package server

import (
	"socket/docs"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/swagger"
	"github.com/swaggo/swag"
)

func (s *Server) initRoutes() {
	s.initSocket()
	s.initAuth()
	s.initSwagger()
	s.initMessage()
}

func (s *Server) initSocket() {
	swag.Register(docs.SwaggerInfo.InfoInstanceName, docs.SwaggerInfo)
	s.app.Get("/ws", websocket.New(s.handler.socketMessageHandler.InitConnection))

}

func (s *Server) initAuth() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.handler.userHandler.Register)
	authRoutes.Post("/login", s.handler.userHandler.Login)
}

func (s *Server) initSwagger() {
	s.app.Get("/swagger/*", swagger.HandlerDefault) // default
}

func (s *Server) initMessage() {
	messageRoutes := s.app.Group("/messages")
	messageRoutes.Get("/", s.handler.socketMessageHandler.GetAll)
}
