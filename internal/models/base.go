package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Base defines a base model for gorm
type Base struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time  `gorm:"index;not null;default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"index;default:null"`
}
