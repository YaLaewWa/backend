package server

import (
	"socket/internal/core/ports"
	"socket/internal/repository"
)

type Repository struct {
	userRepository ports.UserRepository
}

func (s *Server) initRepository() {
	userRepo := repository.NewUserRepo(s.pgDB)
	s.repository = &Repository{
		userRepository: userRepo,
	}
}
