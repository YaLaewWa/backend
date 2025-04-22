package domain

import (
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
