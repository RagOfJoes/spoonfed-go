package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	// ErrFilterInput defines an error where a filter input fails to conver to
	// bson
	ErrFilterInput = errors.New("failed to convert filter input to bson")
)

// See: https://github.com/svett/golang-design-patterns/tree/master/creational-patterns/factory-method
// Follows the design pattern provided above.

// FilterInput defines an interface that stores a collection of
// FilterInput types ie:
// `
//  type RecipeFilterInput struct {
//		UserName *StringFilterInput
//	}
// `
type FilterInput interface {
	Bson() (bson.M, bool)
}

// FilterInputFactory produces FilterInput given an arbritrary
// given argument
type FilterInputFactory interface {
	Create(entries interface{}) FilterInput
}

// FilterPipeline defines a struct that enables building
// match, lookup, and sort aggregation pipelines for MongoDB
type FilterPipeline struct {
	Filters []FilterInput
	// Prerequisites defines a map that takes an arbritrary
	// key from either a Mongo Sort pipeline or a
	// Mongo Match pipeline
	Prerequisites map[string][]bson.M
}

// BuildMatch builds a match aggregation pipeline for MongoDB
func (f *FilterPipeline) BuildMatch(cursor bson.E) (val bson.D) {
	match := bson.D{}
	if cursor.Key != "" {
		match = append(match, cursor)
	}
	for _, f := range f.Filters {
		val, ok := f.Bson()
		if ok {
			for key, val := range val {
				match = append(match, bson.E{Key: key, Value: val})
			}
		}
	}
	return match
}

// BuildPreReq builds arbritraty pipelines that enables future pipelines
// to execute properly
func (f *FilterPipeline) BuildPreReq(sort bson.D, match bson.D) []bson.M {
	prereq := []bson.M{}
	sortMap := sort.Map()
	matchMap := match.Map()
	for key, val := range f.Prerequisites {
		_, okSort := sortMap[key].(int)
		_, okMatch := matchMap[key].(bson.M)
		if okSort || okMatch {
			prereq = append(prereq, val...)
		}
	}
	return prereq
}
