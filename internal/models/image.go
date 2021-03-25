package models

import (
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	URL        string    `gorm:"not null"`
	Name       string    `gorm:"not null;size:72"`
	Caption    string    `gorm:"size:72"`
	EntityType string    `gorm:"not null"`
	EntityID   uuid.UUID `gorm:"index;not null;"`
}

// 								Helpers
// ###########################################

func (i *Image) IsValid(new bool) (ok bool, err error) {
	if !util.CheckLen(i.Name, 1, 72) {
		return false, ErrInvalidLength("image.name", 1, 72)
	} else if !util.IsAlphaNumeric(i.Name, true) {
		return false, ErrInvalidAlphanumeric("image.name")
	}
	if !util.IsURL(i.URL) {
		return false, ErrImageInvalidURL
	}
	return true, nil
}

// 									Hooks
// ###########################################

func (i *Image) BeforeCreate(tx *gorm.DB) (err error) {
	// validate
	if ok, err := i.IsValid(true); !ok {
		return err
	}
	return nil
}
