package server

import (
	"socket/internal/core/ports"
	"socket/internal/handlers"
)

type Handler struct {
	userHandler          ports.UserHandler
	socketMessageHandler ports.MessageSocketHandler
	chatHandler          ports.ChatHandler
}

func (s *Server) initHandler() {
	userHandler := handlers.NewUserHandler(s.service.userService)
	socketMessageHandler := handlers.NewMessageSocketHandler(s.messageHub, s.service.messageService, s.service.chatService)
	chatHandler := handlers.NewChatHandler(s.service.chatService, s.messageHub)
	s.handler = &Handler{
		userHandler:          userHandler,
		socketMessageHandler: socketMessageHandler,
		chatHandler:          chatHandler,
	}
}
