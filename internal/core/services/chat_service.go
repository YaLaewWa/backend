package services

import (
	"errors"
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"

	"github.com/google/uuid"
)

type ChatService struct {
	repo ports.ChatRepository
}

func NewChatService(repo ports.ChatRepository) ports.ChatService {
	return &ChatService{repo: repo}
}

func (c *ChatService) AddUserToChat(chatID uuid.UUID, userID uuid.UUID) (*domain.Chat, error) {
	// Check if chat exist or not
	chat, err := c.repo.GetByID(chatID)
	if err != nil {
		return nil, err
	}

	// Check if the chat is a group chat or not
	if !chat.IsGroup {
		return nil, apperror.ForbiddenError(errors.New("forbidden"), "You are not allowed to join a direct chat")
	}

	// Check if user is already in the chat or not
	isMember, err := c.repo.IsUserInChat(chatID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, apperror.ConflictError(errors.New("conflict"), "You are already in this chat")
	}

	chat, err = c.repo.GetByID(chatID)
	if err != nil {
		return nil, err
	}

	return chat, c.repo.AddUserToChat(chatID, userID)
}

func (c *ChatService) CreateDirectChat(user1 uuid.UUID, user2 uuid.UUID) (*domain.Chat, error) {
	// What name should a direct chat have? If we even have a name for a chat.
	return c.repo.Create("", []uuid.UUID{user1, user2}, false)
}

func (c *ChatService) CreateGroupChat(name string, userIDs []uuid.UUID) (*domain.Chat, error) {
	return c.repo.Create(name, userIDs, true)
}

func (c *ChatService) GetChatsByUserID(userID uuid.UUID, limit int, page int) ([]domain.Chat, int, int, error) {
	return c.repo.GetChatsByUserID(userID, limit, page)
}

func (c *ChatService) GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error) {
	return c.repo.GetChatMembers(chatID, limit, page)
}
