package dto

type UserResponse struct {
	Username string `json:"username"`
}

type UserWithTokenResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"accessToken"`
}

type RegisterRequestBody struct {
	Password string `json:"password" validate:"required,min=8"`
	UserName string `json:"username" validate:"required,min=1"`
}

type LoginRequestBody struct {
	Password string `json:"password" validate:"required"`
	UserName string `json:"username" validate:"required"`
}
