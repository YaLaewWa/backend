package domain

import (
	"socket/internal/dto"
	"time"

	"github.com/google/uuid"
)

type MessageQueue struct {
	Username  string `gorm:"primaryKey" validate:"required"`
	ChatID    uuid.UUID
	Chat      Chat
	Count     int
	UpdatedAt time.Time
}

func (m *MessageQueue) ToDTO() dto.QueueResponse {
	return dto.QueueResponse{
		Username:  m.Username,
		Chat:      m.Chat.ToDTO(),
		Count:     m.Count,
		UpdatedAt: m.UpdatedAt,
	}
}
