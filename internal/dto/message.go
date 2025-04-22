package dto

import (
	"time"

	"github.com/google/uuid"
)

type MessageResponse struct {
	CreateAt time.Time `json:"create_at"`
	Username string    `json:"username"`
	Content  string    `json:"message"`
}

type MessageRequest struct {
	ChatID  uuid.UUID `json:"chat_id"`
	Content string    `json:"content"`
}
