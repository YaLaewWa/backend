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
	s.initChat()
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

func (s *Server) initChat() {
	chatRoutes := s.app.Group("/chats", s.middleware.Auth)
	chatRoutes.Get("/", s.handler.chatHandler.GetChats)
	chatRoutes.Get("/:id/messages", s.handler.socketMessageHandler.GetByChatID)
	chatRoutes.Get("/:id/members", s.handler.chatHandler.GetChatMembers)
	chatRoutes.Post("/direct", s.handler.chatHandler.CreateDirectChat)
	chatRoutes.Post("/group", s.handler.chatHandler.CreateGroupChat)
	chatRoutes.Post("/:id/join", s.handler.chatHandler.JoinChat)
}
