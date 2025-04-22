package ports

import (
	"socket/internal/core/domain"

	"github.com/google/uuid"
)

type MessageQueueRepository interface {
	Create(queue *domain.MessageQueue) error
	Get(chatID uuid.UUID, limit int, page int) ([]domain.MessageQueue, error)
	Update(queue *domain.MessageQueue) error
}
