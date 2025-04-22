package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageQueueRepository struct {
	db *gorm.DB
}

func NewMessageQueueRepository(db *gorm.DB) ports.MessageQueueRepository {
	return &MessageQueueRepository{db: db}
}

func (c *MessageQueueRepository) Create(queue *domain.MessageQueue) error {
	if err := c.db.Create(queue).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to create queue")
	}
	return nil
}

func (c *MessageQueueRepository) GetAll(username string) ([]domain.MessageQueue, error) {
	var queue []domain.MessageQueue
	if err := c.db.Model(queue).Preload("Chat").Preload("Chat.Members").Where("username = ?", username).Find(&queue).Error; err != nil {
		return nil, err
	}
	return queue, nil
}

func (c *MessageQueueRepository) GetQueue(username string, chatID uuid.UUID) (*domain.MessageQueue, error) {
	var queue *domain.MessageQueue

	if err := c.db.Where(&domain.MessageQueue{Username: username, ChatID: chatID}).First(&queue).Error; err != nil {
		return nil, err
	}
	return queue, nil
}

func (c *MessageQueueRepository) Update(queue *domain.MessageQueue) error {
	if err := c.db.Save(&queue).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to update queue")
	}
	return nil
}
