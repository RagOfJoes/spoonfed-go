package models

import (
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

const (
	UsernamMaxLength = 20
)

type User struct {
	Base
	Avatar     *string
	Bio        *string `gorm:"size:72"`
	GivenName  string  `gorm:"not null"`
	FamilyName string  `gorm:"not null"`
	Email      string  `gorm:"not null;uniqueIndex"`
	Username   *string `gorm:"size:20;uniqueIndex"`
	// Set up fields for future so that we can allow
	// for different auth providers
	GithubID    *string `gorm:"uniqueIndex"`
	GoogleID    *string `gorm:"uniqueIndex"`
	RagOfJoesID *string `gorm:"uniqueIndex"`
}

func (o *User) IsValid(new bool) (ok bool, err error) {
	if o == nil {
		return false, ErrUserNotBuilt
	}
	if email := o.Email; !util.IsEmail(email) {
		return false, ErrInvalidEmail
	}
	if avatar := o.Avatar; avatar != nil && !util.IsURL(*avatar) {
		return false, ErrInvalidAvatar
	}
	if !new {
		err := o.validateUsername(*o.Username)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (o *User) validateUsername(username string) error {
	if len(username) < 4 || len(username) > 20 {
		return ErrInvalidLength("user.username", 4, 20)
	}
	if !util.IsAlphaNumeric(username, true) {
		return ErrInvalidAlphanumeric("user.username")
	}
	return nil
}
