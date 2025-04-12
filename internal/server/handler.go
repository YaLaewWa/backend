package server

import (
	"socket/internal/core/ports"
	"socket/internal/handlers"
)

type Handler struct {
	userHandler          ports.UserHandler
	socketMessageHandler ports.MessageSocketHandler
}

func (s *Server) initHandler() {
	userHandler := handlers.NewUserHandler(s.service.userService)
	socketMessageHandler := handlers.NewMessageSocketHandler(s.messageHub)
	s.handler = &Handler{
		userHandler:          userHandler,
		socketMessageHandler: socketMessageHandler,
	}
}
