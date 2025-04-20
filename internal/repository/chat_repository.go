package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ports.ChatRepository {
	return &ChatRepository{db: db}
}

func (c *ChatRepository) Create(name string, users []domain.User, isGroup bool) (*domain.Chat, error) {
	chat := domain.Chat{
		Name:    name,
		IsGroup: isGroup,
		Members: users,
	}
	if err := c.db.Create(&chat).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to create a chat")
	}
	return &chat, nil
}

// func (c *ChatRepository) GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error) {
// 	var chat domain.Chat
// 	if err := c.db.Preload("Members").First(&chat, chatID).Error
// }
