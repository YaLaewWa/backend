package domain

import (
	"socket/internal/dto"
	"time"
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
