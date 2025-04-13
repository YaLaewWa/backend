package services

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
)

type MessageService struct {
	repo ports.MessageRepository
}

func NewMessageRepository(repo ports.MessageRepository) ports.MessageRepository {
	return &MessageService{repo: repo}
}

func (m *MessageService) Create(msg *domain.Message) error {
	return m.repo.Create(msg)
}

func (m *MessageService) GetAll() ([]domain.Message, error) {
	return m.repo.GetAll()
}
