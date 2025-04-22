package domain

import (
	"socket/internal/dto"
	"time"

	"github.com/google/uuid"
)

type User struct {
	CreateAt time.Time `gorm:"autoCreateTime"`
	Username string    `gorm:"primaryKey" validate:"required"`
	Password string    `validate:"required,min=8"`
}

func (u *User) ToDTO() *dto.UserResponse {
	dto := &dto.UserResponse{
		Username: u.Username,
	}
	return dto
}

type Member struct {
	ChatID   uuid.UUID
	Username string
}

func (m *Member) ToDTO() dto.UserResponse {
	return dto.UserResponse{
		Username: m.Username,
	}
}
