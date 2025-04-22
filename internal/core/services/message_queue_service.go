package services

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"

	"github.com/google/uuid"
)

type MessageQueueService struct {
	repo ports.MessageQueueRepository
}

func NewMessageQueueService(repo ports.MessageQueueRepository) ports.MessageQueueService {
	return &MessageQueueService{repo: repo}
}

func (c *MessageQueueService) Create(username string, chatID uuid.UUID) error {
	return nil
}

func (c *MessageQueueService) Get(username string) ([]domain.MessageQueue, error) {
	return nil, nil
}

func (c *MessageQueueService) ReceiveMessage(username string, chatID uuid.UUID) error {
	return nil
}

func (c *MessageQueueService) ReadMessage(username string, chatID uuid.UUID) error {
	return nil
}
