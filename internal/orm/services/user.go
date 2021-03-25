package services

import (
	"context"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"gorm.io/gorm"
)

// 								Dataloaders
// ###########################################

func UserDataloader(ctx context.Context, db *gorm.DB, keys []string) ([]*model.User, []error) {
	users := []models.User{}
	errs := []error{}
	// Execute query
	// then check for Error
	tx := db.Find(&users, "id IN(?)", keys)
	if tx.Error != nil {
		errs = append(errs, tx.Error)
		return nil, errs
	}
	// Map users
	usersArr := []*model.User{}
	usersMap := map[string]*model.User{}
	for _, user := range users {
		bU, err := model.BuildUser(&user)
		if err != nil {
			errs = append(errs, err)
		} else {
			usersMap[user.ID.String()] = bU
		}
	}
	// Append users to array
	// in the order that they
	// were passed
	for _, id := range keys {
		usersArr = append(usersArr, usersMap[id])
	}
	// Return erros if any
	if len(errs) > 0 {
		return nil, errs
	}
	return usersArr, nil
}
