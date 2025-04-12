package server

import (
	"socket/internal/core/ports"
	"socket/internal/core/services"
)

type Service struct {
	userService ports.UserService
}

func (s *Server) initService() {
	userService := services.NewUserService(s.repository.userRepository, s.jwt)
	s.service = &Service{
		userService: userService,
	}
}
