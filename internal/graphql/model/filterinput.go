package model

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrFilterInput = errors.New("failed to convert filter input to bson")
)

// IDFilterInput custom type
type IDFilterInput struct {
	Is     string   `json:"is"`
	NotIs  string   `json:"notIs"`
	Has    []string `json:"has"`
	NotHas []string `json:"notHas"`
}

func (i *IDFilterInput) getCurrentKey() *string {
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		if name == "Is" || name == "NotIs" {
			_, ok := val.FieldByName(name).Interface().(string)
			if ok {
				return &name
			}
		} else if name == "Has" || name == "NotHas" {
			_, ok := val.FieldByName(name).Interface().([]string)
			if ok {
				return &name
			}
		}
	}
	return nil
}

// Key returns the valid field that the key belongs to
func (i *IDFilterInput) Key() string {
	key := i.getCurrentKey()
	switch *key {
	case "NotIs":
		return "$ne"
	case "Has":
		return "$in"
	case "NotHas":
		return "$nin"
	default:
		return "$eq"
	}
}

// Value returns the IDFilterInput value
func (i *IDFilterInput) Value() interface{} {
	ref := reflect.Indirect(reflect.ValueOf(i))
	key := *i.getCurrentKey()
	return ref.FieldByName(key).Interface()
}

// Bson returns the filter interface
// to use with mongo
func (i *IDFilterInput) Bson() (*bson.D, error) {
	key := i.Key()
	if key == "$eq" || key == "$ne" {
		val := i.Value().(string)
		id, _ := primitive.ObjectIDFromHex(val)
		return &bson.D{{
			Key: key, Value: id,
		}}, nil
	} else if key == "$in" || key == "$nin" {
		val := i.Value().([]string)
		ids := []primitive.ObjectID{}
		for _, v := range val {
			id, _ := primitive.ObjectIDFromHex(v)
			ids = append(ids, id)
		}
		return &bson.D{{
			Key: key, Value: ids,
		}}, nil
	}
	return nil, ErrFilterInput
}

// UnmarshalInputIDFilterInput defines custom unmarshal fn
func (i *IDFilterInput) UnmarshalInputIDFilterInput(v interface{}) error {
	return nil
}

type StringFilterInput struct {
	Contains string `json:"contains"`
}

func (s *StringFilterInput) getCurrentKey() *string {
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		val, ok := val.FieldByName(name).Interface().(*StringFilterInput)
		if ok && val != nil {
			return &name
		}
	}
	return nil
}

// Key returns the valid field that the key belongs to
func (s *StringFilterInput) Key() string {
	key := s.getCurrentKey()
	switch *key {
	default:
		return "$eq"
	}
}

// Value returns the IDFilterInput value
func (s *StringFilterInput) Value() string {
	ref := reflect.Indirect(reflect.ValueOf(s))
	key := *s.getCurrentKey()
	return ref.FieldByName(key).Interface().(string)
}

// Bson returns the filter interface
// to use with mongo
func (s *StringFilterInput) Bson() bson.D {
	return bson.D{{
		Key: s.Key(), Value: s.Value(),
	}}
}

// UnmarshalInputIDFilterInput defines custom unmarshal fn
func (i *IDFilterInput) UnmarshalInputStringFilterInput(v interface{}) error {
	return nil
}
