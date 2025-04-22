package dto

import (
	"time"
)

type QueueResponse struct {
	Username  string       `json:"username"`
	Chat      ChatResponse `json:"chat"`
	Count     int          `json:"count"`
	UpdatedAt time.Time    `json:"timestamp"`
}
