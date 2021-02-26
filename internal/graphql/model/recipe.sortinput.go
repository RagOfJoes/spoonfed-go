package model

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// RecipeSortInput defines available
// sort for Recipe queries
type RecipeSortInput struct {
	Name         SortOrder `json:"name"`
	DateCreation SortOrder `json:"dateCreation"`
}

var (
	recipeSortInputKeyMap = SortInputMap{
		"Name": {
			Key:     "name",
			BsonKey: "slug",
		},
		"DateCreation": {
			Key:     "date.creation",
			BsonKey: "date.creation",
		},
	}
)

func (r *RecipeSortInput) getCurrentKey() string {
	val := reflect.ValueOf(r).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		val, ok := val.FieldByName(name).Interface().(SortOrder)
		if ok && val.IsValid() {
			return name
		}
	}
	// Default to DateCreation
	return "DateCreation"
}

func (r *RecipeSortInput) BsonKey() string {
	key := r.getCurrentKey()
	return recipeSortInputKeyMap[key].BsonKey
}

// Key returns the valid field that the key belongs to
func (r *RecipeSortInput) Key() string {
	key := r.getCurrentKey()
	return recipeSortInputKeyMap[key].Key
}

// Value returns the SortOrder value
func (r *RecipeSortInput) Value() SortOrder {
	ref := reflect.Indirect(reflect.ValueOf(r))
	key := r.getCurrentKey()
	return ref.FieldByName(key).Interface().(SortOrder)
}

// Bson returns the sort/filter interface
// to use with mongo
func (r *RecipeSortInput) Bson() bson.D {
	key := r.Key()
	order := r.Value().Int()
	return bson.D{{Key: key, Value: order}}
}

func (r *RecipeSortInput) Cursor(node *Recipe) string {
	key := r.Key()
	switch key {
	case "name":
		return EncodeCursor(node.Slug)
	default:
		unix := int(node.Date.Creation.Unix())
		return EncodeCursor(int(unix))
	}
}

// UnmarshalInputRecipeSortInput defines custom unmarshal fn
func (r *RecipeSortInput) UnmarshalInputRecipeSortInput(v interface{}) error {
	return nil
}
