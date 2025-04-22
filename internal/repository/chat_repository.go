package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/database"
	"socket/pkg/apperror"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ports.ChatRepository {
	return &ChatRepository{db: db}
}

func (c *ChatRepository) Create(name string, userIDs []uuid.UUID, isGroup bool) (*domain.Chat, error) {
	var users []domain.User
	if err := c.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve every user for creating chat")
	}

	chat := domain.Chat{
		Name:    name,
		IsGroup: isGroup,
		Members: users,
	}
	if err := c.db.Create(&chat).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to create chat")
	}
	return &chat, nil
}

func (c *ChatRepository) GetPaginatedChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error) {
	var members []domain.User
	var total, last int

	if err := c.db.Joins("JOIN chat_members cm ON cm.user_id = users.id").Where("cm.chat_id = ?", chatID).
		Scopes(database.Paginate(domain.User{}, &limit, &page, &total, &last)).Find(&members).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve chat's members")
	}
	return members, last, total, nil

}

func (c *ChatRepository) GetAllChatMembers(chatID uuid.UUID) ([]domain.User, error) {
	chat := new(domain.Chat)
	if err := c.db.Preload("Members").First(chat, chatID).Error; err != nil {
		return nil, err
	}
	return chat.Members, nil
}

func (c *ChatRepository) GetAllChatsByUserID(userID uuid.UUID) ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := c.db.Joins("JOIN chat_members cm ON cm.chat_id = chats.id").Where("cm.user_id = ?", userID).Preload("Members").Find(&chats).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve user's chats")
	}
	return chats, nil
}

func (c *ChatRepository) GetPaginatedChatsByUserID(userID uuid.UUID, limit int, page int) ([]domain.Chat, int, int, error) {
	var chats []domain.Chat
	var total, last int

	if err := c.db.Joins("JOIN chat_members cm ON cm.chat_id = chats.id").Where("cm.user_id = ?", userID).
		Preload("Members").Scopes(database.Paginate(&domain.Chat{}, &limit, &page, &total, &last)).
		Find(&chats).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve user's chats")
	}

	return chats, last, total, nil
}

func (c *ChatRepository) AddUserToChat(chatID uuid.UUID, userID uuid.UUID) error {
	chat := domain.Chat{ID: chatID}
	user := domain.User{ID: userID}
	if err := c.db.Model(&chat).Association("Members").Append(&user); err != nil {
		return apperror.InternalServerError(err, "Failed to add user to chat")
	}
	return nil
}

func (c *ChatRepository) GetByID(chatID uuid.UUID) (*domain.Chat, error) {
	chat := new(domain.Chat)
	if err := c.db.Preload("Members").First(chat, chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "Chat not found")
		}
		return nil, apperror.InternalServerError(err, "Failed to retrieve chat")
	}
	return chat, nil
}

func (c *ChatRepository) IsUserInChat(chatID, userID uuid.UUID) (bool, error) {
	var count int64
	err := c.db.Table("chat_members").Where("chat_id = ? AND user_id = ?", chatID, userID).Count(&count).Error
	if err != nil {
		return false, apperror.InternalServerError(err, "Failed to verify membership")
	}
	return count > 0, nil
}
