package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/database"
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

func (m *MessageRepository) GetAll(limit, page int) ([]domain.Message, int, int, error) {
	var msgs []domain.Message
	var total, last int

	if err := m.db.Scopes(database.Paginate(domain.Message{}, &limit, &page, &total, &last)).Order("created_at DESC").Find(&msgs).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve messages")
	}
	return msgs, last, total, nil
}
