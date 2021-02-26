package model

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrInvalidRecipeFilterInputKey = errors.New("invalid recipe filter input key")

	recipeFilterInputKeys = map[string]string{
		"Name":     "name",
		"UserSub":  "createdBy",
		"UserName": "user.username",
	}
)

type RecipeFilterInput struct {
	Name     *StringFilterInput `json:"name"`
	UserSub  *IDFilterInput     `json:"userSub"`
	UserName *StringFilterInput `json:"userName"`
}

type RecipeFilterInputFactory struct{}

func (f *RecipeFilterInputFactory) Create(entries interface{}) FilterInput {
	ent, ok := entries.(RecipeFilterInput)
	if !ok {
		return &RecipeFilterInput{}
	}
	return &ent
}

// Bson returns the filter interface
// to use with mongo
func (r *RecipeFilterInput) Bson() (bson.M, bool) {
	val := reflect.ValueOf(r).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		key := recipeFilterInputKeys[name]
		idf, idfOk := val.FieldByName(name).Interface().(*IDFilterInput)
		strf, strfOk := val.FieldByName(name).Interface().(*StringFilterInput)
		if idfOk && idf != nil {
			return idf.Bson(key)
		} else if strfOk && strf != nil {
			return strf.Bson(key)
		}
	}

	return bson.M{}, false
}
