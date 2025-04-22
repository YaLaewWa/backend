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

	if err = c.repo.AddUserToChat(chatID, userID); err != nil {
		return nil, err
	}

	return c.repo.GetByID(chatID)
}

func (c *ChatService) CreateChat(name string, userIDs []uuid.UUID, isGroup bool) (*domain.Chat, error) {
	return c.repo.Create(name, userIDs, isGroup)
}

func (c *ChatService) GetChatsByUserID(userID uuid.UUID, limit int, page int) ([]domain.Chat, int, int, error) {
	if limit <= 0 {
		chats, err := c.repo.GetAllChatsByUserID(userID)
		return chats, 1, len(chats), err
	}
	return c.repo.GetPaginatedChatsByUserID(userID, limit, page)
}

func (c *ChatService) GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error) {
	if limit <= 0 {
		members, err := c.repo.GetAllChatMembers(chatID)
		return members, 1, len(members), err
	}
	return c.repo.GetPaginatedChatMembers(chatID, limit, page)
}
