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
}

func (s *Server) initSocket() {
	s.app.Use("/ws", s.middleware.Auth, s.middleware.Websocket)
	s.app.Get("/ws", websocket.New(s.handler.socketMessageHandler.InitConnection))

}

func (s *Server) initAuth() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.handler.userHandler.Register)
	authRoutes.Post("/login", s.handler.userHandler.Login)
}

func (s *Server) initSwagger() {
	swag.Register(docs.SwaggerInfo.InfoInstanceName, docs.SwaggerInfo)
	s.app.Get("/swagger/*", swagger.HandlerDefault) // default
}
