package model

import (
	"fmt"

	"github.com/RagOfJoes/spoonfed-go/internal/models"
)

type User struct {
	ID       string    `json:"id"`
	Username *string   `json:"username"`
	Bio      *string   `json:"bio"`
	Avatar   *string   `json:"avatar"`
	Email    Email     `json:"email"`
	Name     UserName  `json:"name"`
	Date     *MetaDate `json:"date"`
}

func BuildUser(u *models.User) (*User, error) {
	ok, err := u.IsValid(false)
	if !ok || err != nil {
		return nil, err
	}
	full := fmt.Sprintf("%s %s", u.GivenName, u.FamilyName)
	user := &User{
		ID:       u.ID.String(),
		Bio:      u.Bio,
		Avatar:   u.Avatar,
		Username: u.Username,
		Email:    Email(u.Email),
		Date: &MetaDate{
			Creation:   u.CreatedAt,
			LastUpdate: u.UpdatedAt,
		},
		Name: UserName{
			First: &u.GivenName,
			Last:  &u.FamilyName,
			Full:  &full,
		},
	}
	return user, nil
}
