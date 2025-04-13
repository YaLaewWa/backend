package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	UpdateAt time.Time `gorm:"autoUpdateTime"`
	Username string
	Content  string
	// TODO: connect message to some chat
}
