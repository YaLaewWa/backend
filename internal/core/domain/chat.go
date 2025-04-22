package domain

import (
	"socket/internal/dto"

	"github.com/google/uuid"
)

type Chat struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name    string
	IsGroup bool
	Members []User `gorm:"many2many:chat_members;"`
}

func (c *Chat) ToDTO() dto.ChatResponse {
	members := make([]dto.UserResponse, len(c.Members))
	for i, member := range c.Members {
		members[i] = *member.ToDTO()
	}

	return dto.ChatResponse{
		ID:      c.ID,
		Name:    c.Name,
		IsGroup: c.IsGroup,
		Members: members,
		Joined:  true,
	}
}

type ChatWithMembership struct {
	ID      uuid.UUID
	Name    string
	IsGroup bool
	Joined  bool
}

func (c *ChatWithMembership) ToDTO(members []dto.UserResponse) dto.ChatResponse {
	return dto.ChatResponse{
		ID:      c.ID,
		Name:    c.Name,
		IsGroup: c.IsGroup,
		Joined:  c.Joined,
		Members: members,
	}
}

func (c *Chat) ToSocketDTO(isCreater bool) dto.ChatSocket {
	members := make([]dto.UserResponse, len(c.Members))
	for i, member := range c.Members {
		members[i] = *member.ToDTO()
	}

	return dto.ChatSocket{
		Type: "new_group",
		Payload: dto.ChatPayload{
			ID:      c.ID,
			Name:    c.Name,
			Joined:  isCreater,
			Members: members,
		},
	}
}
