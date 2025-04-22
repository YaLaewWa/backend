package server

import (
	"socket/internal/core/ports"
	"socket/internal/handlers"
)

type Handler struct {
	userHandler          ports.UserHandler
	socketMessageHandler ports.MessageSocketHandler
	chatHandler          ports.ChatHandler
	messageQueueHandler  ports.MessageQueueHandler
}

func (s *Server) initHandler() {
	userHandler := handlers.NewUserHandler(s.service.userService)
	socketMessageHandler := handlers.NewMessageSocketHandler(s.messageHub, s.service.messageService, s.service.chatService)
	messageQueueHandler := handlers.NewMessageQueueHandler(s.service.messageQueueService)
	chatHandler := handlers.NewChatHandler(s.service.chatService, s.service.messageQueueService, s.messageHub)
	s.handler = &Handler{
		userHandler:          userHandler,
		socketMessageHandler: socketMessageHandler,
		chatHandler:          chatHandler,
		messageQueueHandler:  messageQueueHandler,
	}
}
