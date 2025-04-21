package dto

import "github.com/google/uuid"

type ChatResponse struct {
	ID      uuid.UUID      `json:"id"`
	Name    string         `json:"name,omitempty"`
	IsGroup bool           `json:"is_group"`
	Members []UserResponse `json:"members"`
}

type CreateChatRequest struct {
	Name    string      `json:"name"`
	UserIDs []uuid.UUID `json:"user_ids" validate:"required,min=1"`
}
