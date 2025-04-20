package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/database"
	"socket/pkg/apperror"

	"github.com/google/uuid"
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

func (c *ChatRepository) GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error) {
	var members []domain.User
	var total, last int

	if err := c.db.Joins("JOIN chat_members cm ON cm.user_id = users.id").Where("cm.chat_id = ?", chatID).
		Scopes(database.Paginate(domain.User{}, &limit, &page, &total, &last)).Find(&members).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve chat's members")
	}
	return members, last, total, nil

}

func (c *ChatRepository) GetChatByUserID(userID uuid.UUID, limit int, page int) ([]domain.Chat, int, int, error) {
	var chats []domain.Chat
	var total, last int

	if err := c.db.Joins("JOIN chat_members cm ON cm.chat_id = chats.id").Where("cm.user_id = ?", userID).
		Preload("Members").Scopes(database.Paginate(&domain.Chat{}, &limit, &page, &total, &last)).
		Find(&chats).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve user's chats")
	}

	return chats, last, total, nil
}
