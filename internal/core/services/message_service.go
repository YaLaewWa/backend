package services

import (
	"errors"
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/pkg/apperror"

	"github.com/google/uuid"
)

type MessageService struct {
	msgRepo  ports.MessageRepository
	chatRepo ports.ChatRepository
}

func NewMessageService(msgRepo ports.MessageRepository, chatRepo ports.ChatRepository) ports.MessageService {
	return &MessageService{msgRepo: msgRepo, chatRepo: chatRepo}
}

func (m *MessageService) Create(sender string, chatID uuid.UUID, content string) (*domain.Message, error) {
	msg := &domain.Message{
		Username: sender,
		ChatID:   chatID,
		Content:  content,
	}
	return msg, m.msgRepo.Create(msg)
}

func (m *MessageService) GetByChatID(chatID uuid.UUID, limit int, page int, username string) ([]domain.Message, int, int, error) {
	// Check if chat exist or not
	_, err := m.chatRepo.GetByID(chatID)
	if err != nil {
		return nil, 0, 0, err
	}

	// Check if user is a member of chat or not
	isMember, err := m.chatRepo.IsUserInChat(chatID, username)
	if err != nil {
		return nil, 0, 0, err
	}
	if !isMember {
		return nil, 0, 0, apperror.ForbiddenError(errors.New("forbidden"), "You are not a member of this chat")
	}

	if limit <= 0 {
		msgs, err := m.msgRepo.GetAllByChatID(chatID)
		return msgs, 1, len(msgs), err
	}

	return m.msgRepo.GetPaginatedByChatID(chatID, limit, page)
}
