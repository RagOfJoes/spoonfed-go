package model

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// CursorSortInput defines custom CursorSortinput
// model
type CursorSortInput struct {
	Name         *SortOrder `json:"name"`
	DateCreation *SortOrder `json:"dateCreation"`
}

func (c *CursorSortInput) getCurrentKey() *string {
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		val, ok := val.FieldByName(name).Interface().(*SortOrder)
		if ok && val != nil {
			return &name
		}
	}
	return nil
}

// Key returns the valid field that the key belongs to
func (c *CursorSortInput) Key() string {
	key := c.getCurrentKey()
	switch *key {
	case "Name":
		return "name"
	default:
		return "date.creation"
	}
}

// Value returns the SortOrder value
func (c *CursorSortInput) Value() *SortOrder {
	ref := reflect.Indirect(reflect.ValueOf(c))
	key := *c.getCurrentKey()
	return ref.FieldByName(key).Interface().(*SortOrder)
}

// Bson returns the sort/filter interface
// to use with mongo
func (c *CursorSortInput) Bson() bson.D {
	return bson.D{{
		Key: c.Key(), Value: c.Value().Int(),
	}}
}

// UnmarshalInputCursorSortInput defines custom unmarshal fn
func (c *CursorSortInput) UnmarshalInputCursorSortInput(v interface{}) error {
	return nil
}
