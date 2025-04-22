package ports

import (
	"socket/internal/core/domain"

	"github.com/google/uuid"
)

type MessageQueueRepository interface {
	Create(queue *domain.MessageQueue) error
	GetAll(username string) ([]domain.MessageQueue, error)
	GetQueue(username string, chatID uuid.UUID) (*domain.MessageQueue, error)
	Update(queue *domain.MessageQueue) error
}

type MessageQueueService interface {
	Create(username string, chatID uuid.UUID) error
	Get(username string) ([]domain.MessageQueue, error)
	ReceiveMessage(username string, chatID uuid.UUID) error
	ReadMessage(username string, chatID uuid.UUID) error
}
