package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) ports.MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) Create(msg *domain.Message) error {
	if err := m.db.Create(msg).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to save message to database")
	}
	return nil
}

func (m *MessageRepository) GetAll() ([]domain.Message, error) {
	var msgs []domain.Message

	// Not to sure what the order suppose to be here
	// Also does chat have pagination?
	if err := m.db.Order("created_at DESC").Find(&msgs).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve messages")
	}
	return msgs, nil
}
