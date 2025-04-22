package ports

import (
	"socket/internal/core/domain"
)

type MessageQueueRepository interface {
	Create(queue *domain.MessageQueue) error
	Get(username string) ([]domain.MessageQueue, error)
	Update(queue *domain.MessageQueue) error
}
