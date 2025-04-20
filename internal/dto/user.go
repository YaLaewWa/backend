package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type UserWithTokenResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"accessToken"`
}

type AuthBody struct {
	Password string `json:"password" validate:"required,min=9"`
	UserName string `json:"username" validate:"required,min=4"`
}
