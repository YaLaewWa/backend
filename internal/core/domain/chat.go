package domain

import (
	"socket/internal/dto"

	"github.com/google/uuid"
)

type Chat struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name    string    // Is this field needed? I just thought of this and it isn't in the design.
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
	}
}
