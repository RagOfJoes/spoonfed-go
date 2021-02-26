package model

import "go.mongodb.org/mongo-driver/bson"

// SortInputKey defines a map value
// that allows for faster lookup
// times for Key, BsonKey, Value fns
type SortInputKey struct {
	Key            string
	BsonKey        string
	AdditionalSort func(sortKey string, sortOrder int) []bson.E
}

type SortInputMap map[string]SortInputKey
