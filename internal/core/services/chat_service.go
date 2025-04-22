package services

import (
	"errors"
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/dto"
	"socket/pkg/apperror"

	"github.com/google/uuid"
)

type ChatService struct {
	repo ports.ChatRepository
}

func NewChatService(repo ports.ChatRepository) ports.ChatService {
	return &ChatService{repo: repo}
}

func (c *ChatService) AddUserToChat(chatID uuid.UUID, username string) (*domain.Chat, error) {
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
	isMember, err := c.repo.IsUserInChat(chatID, username)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, apperror.ConflictError(errors.New("conflict"), "You are already in this chat")
	}

	if err = c.repo.AddUserToChat(chatID, username); err != nil {
		return nil, err
	}

	return c.repo.GetByID(chatID)
}

func (c *ChatService) CreateChat(name string, usernames []string, isGroup bool) (*domain.Chat, error) {
	if !isGroup {
		dm, err := c.repo.GetDMChat(usernames)
		if err != nil {
			return nil, err
		}
		if dm != nil {
			return nil, apperror.ConflictError(errors.New("conflict"), "You already have a direct message chat with this user")
		}
	}
	return c.repo.Create(name, usernames, isGroup)
}

func (c *ChatService) GetChatsByUsername(username string, limit int, page int) ([]domain.Chat, int, int, error) {
	if limit <= 0 {
		chats, err := c.repo.GetAllChatsByUsername(username)
		return chats, 1, len(chats), err
	}
	return c.repo.GetPaginatedChatsByUsername(username, limit, page)
}

func (c *ChatService) GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error) {
	if limit <= 0 {
		members, err := c.repo.GetAllChatMembers(chatID)
		return members, 1, len(members), err
	}
	return c.repo.GetPaginatedChatMembers(chatID, limit, page)
}

func (c *ChatService) IsUserInChat(chatID uuid.UUID, username string) (bool, error) {
	return c.repo.IsUserInChat(chatID, username)
}

func (c *ChatService) GetGroupChats(username string, limit int, page int) ([]dto.ChatResponse, int, int, error) {
	if limit <= 0 {
		chats, err := c.repo.GetAllGroupChats(username)
		return chats, 1, len(chats), err
	}
	return c.repo.GetPaginatedGroupChats(username, limit, page)
}

func (c *ChatService) HavePrivateChat(user1, user2 string) (bool, error) {
	chat1, err := c.repo.GetAllChatsByUsername(user1)
	if err != nil {
		return false, err
	}
	chat2, err := c.repo.GetAllChatsByUsername(user2)
	if err != nil {
		return false, err
	}
	privateChatsUser2 := make(map[uuid.UUID]bool)
	for _, item := range chat2 {
		if !item.IsGroup {
			privateChatsUser2[item.ID] = true
		}
	}
	for _, item := range chat1 {
		if !item.IsGroup {
			if privateChatsUser2[item.ID] {
				return true, nil
			}
		}
	}
	return false, nil
}
