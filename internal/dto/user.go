package dto

type UserResponse struct {
	Username string `json:"username"`
}

type UserWithTokenResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"accessToken"`
}

type RegisterRequestBody struct {
	Password string `json:"password" validate:"required"`
	UserName string `json:"username" validate:"required"`
}

type LoginRequestBody struct {
	Password string `json:"password" validate:"required"`
	UserName string `json:"username" validate:"required"`
}
