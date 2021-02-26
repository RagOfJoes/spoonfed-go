package model

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// StringFilterInput defines a custom type
type StringFilterInput struct {
	Matches  *string `json:"matches"`
	Contains *string `json:"contains"`
}

func (s *StringFilterInput) getField() (key string, value string, ok bool) {
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		v, ok := val.FieldByName(name).Interface().(*string)
		if ok && v != nil {
			return name, *v, true
		}
	}
	return "", "", false
}

// Key returns the valid field that the key belongs to
func (s *StringFilterInput) Key() (key string, ok bool) {
	k, _, o := s.getField()
	if !o {
		return "", false
	}
	switch k {
	case "Contains":
		return "$text", true
	default:
		return "$eq", true
	}
}

// Bson returns the filter interface
// to use with mongo
func (s *StringFilterInput) Bson(key string) (val bson.M, ok bool) {
	bsonKey, ok := s.Key()
	if !ok {
		return bson.M{}, false
	}
	_, value, ok := s.getField()
	if !ok {
		return bson.M{}, false
	}
	if bsonKey == "$text" {
		return bson.M{
			bsonKey: bson.M{
				"$search": value,
			},
		}, true
	}
	return bson.M{key: bson.M{
		bsonKey: value,
	}}, true
}

// UnmarshalInputStringFilterInput defines custom unmarshal fn
func (s *StringFilterInput) UnmarshalInputStringFilterInput(v interface{}) error {
	return nil
}
