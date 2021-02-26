package model

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDFilterInput custom type
type IDFilterInput struct {
	Is     *string   `json:"is"`
	NotIs  *string   `json:"notIs"`
	Has    []*string `json:"has"`
	NotHas []*string `json:"notHas"`
}

func (i *IDFilterInput) getField() (key string, value interface{}, ok bool) {
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		if name == "Is" || name == "NotIs" {
			v, ok := val.FieldByName(name).Interface().(*string)
			if ok && v != nil {
				return name, v, true
			}
		} else if name == "Has" || name == "NotHas" {
			v, ok := val.FieldByName(name).Interface().([]*string)
			if ok && v != nil {
				return name, v, true
			}
		}
	}
	return "", nil, false
}

// Key returns the valid field that the key belongs to
func (i *IDFilterInput) Key() (key string, ok bool) {
	k, _, o := i.getField()
	if !o {
		return k, o
	}

	switch k {
	case "NotIs":
		return "$ne", true
	case "Has":
		return "$in", true
	case "NotHas":
		return "$nin", true
	default:
		return "$eq", true
	}
}

// Bson returns the filter interface
// to use with mongo
func (i *IDFilterInput) Bson(key string) (val bson.M, ok bool) {
	innerKey, ok := i.Key()
	if !ok {
		return bson.M{}, false
	}
	_, intf, ok := i.getField()
	if !ok {
		return bson.M{}, false
	}
	if key == "$eq" || key == "$ne" {
		val := intf.(*string)
		id, _ := primitive.ObjectIDFromHex(*val)
		return bson.M{
			key: bson.M{
				innerKey: id,
			},
		}, true
	} else if key == "$in" || key == "$nin" {
		val := intf.([]*string)
		ids := []primitive.ObjectID{}
		for _, v := range val {
			id, _ := primitive.ObjectIDFromHex(*v)
			ids = append(ids, id)
		}
		return bson.M{
			key: bson.M{
				innerKey: ids,
			},
		}, true
	}
	return bson.M{}, false
}

// UnmarshalInputIDFilterInput defines custom unmarshal fn
func (i *IDFilterInput) UnmarshalInputIDFilterInput(v interface{}) error {
	return nil
}
