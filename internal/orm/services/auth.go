package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"github.com/RagOfJoes/spoonfed-go/pkg/auth"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"gorm.io/gorm"
)

var (
	ErrInvalidIDProvided        = errors.New("invalid id provided")
	ErrInvalidEmailProvided     = errors.New("invalid email provided")
	ErrSignupProcessNotFinished = errors.New("sign up")
)

// IsAuthenticated retrieves `models.User` from context
// and if exists then builds a `model.User` from it to return.
// If it doesn't exist then return err or nil
func IsAuthenticated(ctx context.Context) (*models.User, error) {
	user, err := auth.GetUserFromContext(ctx, util.ProjectContextKeys.User)
	if err != nil {
		return nil, err
	}
	u, ok := user.(*models.User)
	if !ok {
		return nil, nil
	}
	return u, nil
}

// GetUserFromDB does just that
func GetUserFromDB(db *gorm.DB, id string) (*model.User, error) {
	u := &models.User{}
	db = db.First(u, "id = ?", id)
	if db.Error != nil {
		return nil, db.Error
	}
	user, err := model.BuildUser(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserFromRagOfJoes does just that
func GetUserFromRagOfJoes(u *auth.User, o *gorm.DB) (*models.User, error) {
	sub := getSub(u)
	if u == nil || sub == "" {
		return nil, ErrInvalidIDProvided
	}
	email := u.Email
	if email == "" {
		return nil, ErrInvalidEmailProvided
	}
	found := &models.User{}
	first := o.First(found, models.User{RagOfJoesID: &sub})
	err := first.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		placeholderUsername := strings.ToLower(fmt.Sprintf("%s-%s", u.Profile.FirstName, u.Profile.LastName))
		placeholderUsername = util.SmartTruncate(placeholderUsername, 10)
		placeholderUsername += util.SlugID(4)
		user := &models.User{
			RagOfJoesID: &sub,
			Email:       u.Email,
			GivenName:   u.Profile.FirstName,
			FamilyName:  u.Profile.LastName,
			Username:    &placeholderUsername,
		}
		create := o.Create(user)
		if err := create.Error; err != nil {
			return nil, err
		}
		return user, nil
	} else if err != nil {
		return nil, err
	}

	return found, nil
}

func getSub(u *auth.User) string {
	if u.Sub == "" {
		return u.Profile.Sub
	}
	return u.Sub
}
