package dto

import "github.com/google/uuid"

type ChatResponse struct {
	ID      uuid.UUID      `json:"id"`
	Name    string         `json:"name"`
	IsGroup bool           `json:"is_group"`
	Members []UserResponse `json:"members"`
}

type CreateDirectChatRequest struct {
	User1 uuid.UUID `json:"user1_id" validate:"required"`
	User2 uuid.UUID `json:"user2_id" validate:"required"`
}

type CreateGroupChatRequest struct {
	Name    string      `json:"name" validate:"required"`
	UserIDs []uuid.UUID `json:"user_ids" validate:"required,min=2"`
	// technically line/discord can create 1 person group chat but what should be our min?
}
