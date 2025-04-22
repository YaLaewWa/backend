package server

import (
	"socket/internal/core/ports"
	"socket/internal/core/services"
)

type Service struct {
	userService         ports.UserService
	messageService      ports.MessageService
	chatService         ports.ChatService
	messageQueueService ports.MessageQueueService
}

func (s *Server) initService() {
	userService := services.NewUserService(s.repository.userRepository, s.jwt)
	messageService := services.NewMessageService(s.repository.messageRepository, s.repository.chatRepository)
	chatService := services.NewChatService(s.repository.chatRepository)
	messageQueueService := services.NewMessageQueueService(s.repository.messageQueueRepository)
	s.service = &Service{
		userService:         userService,
		messageService:      messageService,
		chatService:         chatService,
		messageQueueService: messageQueueService,
	}
}
