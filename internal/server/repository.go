package server

import (
	"socket/internal/core/ports"
	"socket/internal/repository"
)

type Repository struct {
	userRepository    ports.UserRepository
	messageRepository ports.MessageRepository
	chatRepository    ports.ChatRepository
}

func (s *Server) initRepository() {
	userRepo := repository.NewUserRepo(s.pgDB)
	messageRepo := repository.NewMessageRepository(s.pgDB)
	chatRepo := repository.NewChatRepository(s.pgDB)
	s.repository = &Repository{
		userRepository:    userRepo,
		messageRepository: messageRepo,
		chatRepository:    chatRepo,
	}
}
