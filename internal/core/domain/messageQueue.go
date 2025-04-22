package domain

import (
	"socket/internal/dto"
	"time"

	"github.com/google/uuid"
)

type MessageQueue struct {
	Username  string    `gorm:"primaryKey" validate:"required"`
	ChatID    uuid.UUID `gorm:"primaryKey" validate:"required"`
	Chat      Chat
	Count     int
	Timestamp time.Time
}

func (m *MessageQueue) ToDTO() dto.QueueResponse {
	name := m.Chat.Name
	if !m.Chat.IsGroup {
		if m.Username == m.Chat.Members[0].Username {
			name = m.Chat.Members[1].Username
		} else {
			name = m.Chat.Members[0].Username
		}
	}
	dtoChat := m.Chat.ToDTO()
	dtoChat.Name = name
	return dto.QueueResponse{
		Username:  m.Username,
		Chat:      dtoChat,
		Count:     m.Count,
		UpdatedAt: m.Timestamp,
	}
}
