package dto

import "github.com/google/uuid"

type ChatResponse struct {
	ID      uuid.UUID      `json:"id"`
	Name    string         `json:"name,omitempty"`
	IsGroup bool           `json:"is_group"`
	Joined  bool           `json:"joined"`
	Members []UserResponse `json:"members"`
}

type CreateChatRequest struct {
	Name      string   `json:"name"`
	Usernames []string `json:"usernames" validate:"required,min=1"`
}
