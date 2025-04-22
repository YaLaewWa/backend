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

type ChatSocket struct {
	Type    string      `json:"type"`
	Payload ChatPayload `json:"payload"`
}

type ChatPayload struct {
	ID      uuid.UUID      `json:"id"`
	Name    string         `json:"name"`
	Joined  bool           `json:"joined"`
	Members []UserResponse `json:"members"`
}

type JoinGroupResponse struct {
	Type    string           `json:"type"`
	Payload JoinGroupPayload `json:"payload"`
}

type JoinGroupPayload struct {
	ChatID   uuid.UUID `json:"chat_id"`
	Username string    `json:"username"`
}

func GetJoinSocketDTO(chatID uuid.UUID, username string) JoinGroupResponse {
	return JoinGroupResponse{
		Type: "new_user_group",
		Payload: JoinGroupPayload{
			ChatID:   chatID,
			Username: username,
		},
	}
}
