package domain

import "github.com/google/uuid"

type Sample struct {
	Id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
}
