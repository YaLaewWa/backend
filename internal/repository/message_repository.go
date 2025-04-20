package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/database"
	"socket/pkg/apperror"

	"github.com/google/uuid"
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

func (m *MessageRepository) GetByChatID(chatID uuid.UUID, limit int, page int) ([]domain.Message, int, int, error) {
	var msgs []domain.Message
	var total, last int

	if err := m.db.Where("chat_id = ?", chatID).Order("create_at DESC").
		Scopes(database.Paginate(domain.Message{}, &limit, &page, &total, &last)).Find(&msgs).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve messages")
	}

	return msgs, last, total, nil
}
