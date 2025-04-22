package dto

import (
	"time"

	"github.com/google/uuid"
)

type MessageResponse struct {
	Type    string                 `json:"type"`
	Payload MessageResponsePayload `json:"payload"`
}

type MessageResponsePayload struct {
	CreateAt time.Time `json:"create_at"`
	Username string    `json:"username"`
	Content  string    `json:"message"`
	ChatID   uuid.UUID `json:"chat_id"`
}

type MessageRequest struct {
	Type    string                `json:"type"`
	Payload MessageRequestPayload `json:"payload"`
}

type MessageRequestPayload struct {
	ChatID  uuid.UUID `json:"chat_id"`
	Content string    `json:"content"`
}
