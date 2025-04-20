package ports

import (
	"socket/internal/core/domain"
)

type ChatRepository interface {
	Create(name string, users []domain.User, isGroup bool) (*domain.Chat, error)
	// GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error)
}
