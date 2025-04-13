package dto

import "time"

type MessageResponse struct {
	CreateAt time.Time `json:"create_at"`
	Username string    `json:"username"`
	Content  string    `json:"message"`
}
