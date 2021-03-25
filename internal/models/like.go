package models

import "github.com/gofrs/uuid"

type Like struct {
	Base
	Active     int8      `gorm:"default:1"`
	EntityType string    `gorm:"not null"`
	EntityID   uuid.UUID `gorm:"index;not null;"`
	UserID     uuid.UUID `gorm:"index;not null;"`
}
