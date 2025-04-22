package repository

import (
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/database"
	"socket/internal/dto"
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

func (c *ChatRepository) Create(name string, usernames []string, isGroup bool) (*domain.Chat, error) {
	var users []domain.User
	if err := c.db.Where("username IN ?", usernames).Find(&users).Error; err != nil {
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

	if err := c.db.Joins("JOIN chat_members cm ON cm.user_username = users.username").Where("cm.chat_id = ?", chatID).
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

func (c *ChatRepository) GetAllChatsByUsername(username string) ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := c.db.Joins("JOIN chat_members cm ON cm.chat_id = chats.id").Where("cm.user_username = ?", username).Preload("Members").Find(&chats).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve user's chats")
	}
	return chats, nil
}

func (c *ChatRepository) GetPaginatedChatsByUsername(username string, limit int, page int) ([]domain.Chat, int, int, error) {
	var chats []domain.Chat
	var total, last int

	if err := c.db.Joins("JOIN chat_members cm ON cm.chat_id = chats.id").Where("cm.user_username = ?", username).
		Preload("Members").Scopes(database.Paginate(&domain.Chat{}, &limit, &page, &total, &last)).
		Find(&chats).Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve user's chats")
	}

	return chats, last, total, nil
}

func (c *ChatRepository) AddUserToChat(chatID uuid.UUID, username string) error {
	chat := domain.Chat{ID: chatID}
	user := domain.User{Username: username}
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

func (c *ChatRepository) IsUserInChat(chatID uuid.UUID, username string) (bool, error) {
	var count int64
	err := c.db.Table("chat_members").Where("chat_id = ? AND user_username = ?", chatID, username).Count(&count).Error
	if err != nil {
		return false, apperror.InternalServerError(err, "Failed to verify membership")
	}
	return count > 0, nil
}

func (c *ChatRepository) GetPaginatedGroupChats(username string, limit int, page int) ([]dto.ChatResponse, int, int, error) {
	var res []domain.ChatWithMembership
	var total, last int

	err := c.db.Table("chats c").Select("c.id, c.name, c.is_group, CASE WHEN cm.user_username IS NULL THEN false ELSE true END as joined").
		Joins("LEFT JOIN chat_members cm ON cm.chat_id = c.id AND cm.user_username = ?", username).
		Where("c.is_group = ?", true).
		Scopes(database.Paginate(&domain.ChatWithMembership{}, &limit, &page, &total, &last)).Scan(&res).Error
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve group chats")
	}

	groups, err := c.preloadMembers(res)
	return groups, last, total, err
}

func (c *ChatRepository) GetAllGroupChats(username string) ([]dto.ChatResponse, error) {
	var res []domain.ChatWithMembership

	err := c.db.Table("chats c").Select("c.id, c.name, c.is_group, CASE WHEN cm.user_username IS NULL THEN false ELSE true END as joined").
		Joins("LEFT JOIN chat_members cm ON cm.chat_id = c.id AND cm.user_username = ?", username).
		Where("c.is_group = ?", true).Scan(&res).Error
	if err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve group chats")
	}

	return c.preloadMembers(res)
}

func (c *ChatRepository) preloadMembers(res []domain.ChatWithMembership) ([]dto.ChatResponse, error) {
	chatIDs := make([]uuid.UUID, len(res))
	for i, g := range res {
		chatIDs[i] = g.ID
	}

	var members []domain.Member
	err := c.db.Table("chat_members").Select("chat_id, users.username").
		Joins("JOIN users ON users.username = chat_members.user_username").
		Where("chat_id IN ?", chatIDs).Scan(&members).Error
	if err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve group chats' members")
	}

	groupMap := make(map[uuid.UUID][]dto.UserResponse)
	for _, m := range members {
		groupMap[m.ChatID] = append(groupMap[m.ChatID], m.ToDTO())
	}

	groups := make([]dto.ChatResponse, len(res))
	for i, group := range res {
		groups[i] = group.ToDTO(groupMap[group.ID])
	}

	return groups, nil
}

func (c *ChatRepository) GetDMChat(usernames []string) (*domain.Chat, error) {
	chat := new(domain.Chat)

	query := c.db.Table("chat_members").Select("chat_id").Where("user_username IN ?", usernames).
		Group("chat_id").Having("COUNT(DISTINCT user_username) = 2")

	err := c.db.Where("id IN (?) AND is_group = false", query).Preload("Members").First(chat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, apperror.InternalServerError(err, "Failed to retrieve chat.")
		}
	}

	return chat, nil
}
