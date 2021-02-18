package model

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrInvalidRecipeFilterInputKey = errors.New("invalid recipe filter input key")
)

type RecipeFilterInput struct {
	UserSub IDFilterInput `json:"userSub"`
}

func (r *RecipeFilterInput) getCurrentKey() *string {
	val := reflect.ValueOf(r).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		switch name {
		case "UserSub":
			_, ok := val.FieldByName(name).Interface().(IDFilterInput)
			if ok {
				return &name
			}
		case "UserName":
			_, ok := val.FieldByName(name).Interface().(StringFilterInput)
			if ok {
				return &name
			}
		}
	}
	return nil
}

// Bson returns the filter interface
// to use with mongo
func (r *RecipeFilterInput) Bson() (*bson.E, error) {
	key := r.getCurrentKey()
	if key == nil {
		return nil, ErrInvalidRecipeFilterInputKey
	}
	switch *key {
	case "UserSub":
		val, err := r.UserSub.Bson()
		if err != nil {
			return nil, err
		}
		return &bson.E{
			Key: "createdBy", Value: val,
		}, nil
	}
	return nil, ErrInvalidRecipeFilterInputKey
}
