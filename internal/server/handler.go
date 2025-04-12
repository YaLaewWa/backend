package server

import (
	"socket/internal/core/ports"
	"socket/internal/handlers"
)

type Handler struct {
	userHandler ports.UserHandler
}

func (s *Server) initHandler() {
	userHandler := handlers.NewUserHandler(s.service.userService)
	s.handler = &Handler{
		userHandler: userHandler,
	}
}
