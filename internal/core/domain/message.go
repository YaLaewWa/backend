package domain

import (
	"socket/internal/dto"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	UpdateAt time.Time `gorm:"autoUpdateTime"`
	Username string
	Content  string
	ChatID   uuid.UUID
}

func (m *Message) ToDTO() dto.MessageResponse {
	return dto.MessageResponse{
		Type:     "message",
		CreateAt: m.CreateAt,
		Username: m.Username,
		Content:  m.Content,
	}
}

type HubMessage struct {
	Message Message
	To      []User
}
