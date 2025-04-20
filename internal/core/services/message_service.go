package services

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"

	"github.com/google/uuid"
)

type MessageService struct {
	repo ports.MessageRepository
}

func NewMessageService(repo ports.MessageRepository) ports.MessageService {
	return &MessageService{repo: repo}
}

func (m *MessageService) Create(msg *domain.Message) error {
	return m.repo.Create(msg)
}

func (m *MessageService) GetByChatID(chatID uuid.UUID, limit int, page int) ([]domain.Message, int, int, error) {
	return m.repo.GetByChatID(chatID, limit, page)
}
