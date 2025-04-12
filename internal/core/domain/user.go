package domain

import (
	"socket/internal/dto"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	UpdateAt time.Time `gorm:"autoUpdateTime"`
	Username string    `gorm:"unique" validate:"required"`
	Password string    `validate:"required,min=8"`
}

func (u *User) ToDTO() *dto.UserResponse {
	dto := &dto.UserResponse{
		Username: u.Username,
	}
	return dto
}
