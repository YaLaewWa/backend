package server

import (
	"socket/internal/core/ports"
	"socket/internal/core/services"
)

type Service struct {
	userService    ports.UserService
	messageService ports.MessageService
}

func (s *Server) initService() {
	userService := services.NewUserService(s.repository.userRepository, s.jwt)
	messageService := services.NewMessageService(s.repository.messageRepository)
	s.service = &Service{
		userService:    userService,
		messageService: messageService,
	}
}
