package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"gorm.io/gorm"
)

var (
	errInvalidFilterCondition = func(key string, conditions []model.FilterCondition) error {
		var toStrings []string
		for _, vc := range conditions {
			toStrings = append(toStrings, vc.String())
		}
		return fmt.Errorf("invalid filter condition provided for key `%s`. valid conditions are: %s", key, strings.Join(toStrings, ", "))
	}
)

type FilterRules struct {
	// Keys define valid FilterInput keys
	Keys []string
	// Rules define custom rule for a specific
	// FilterInput key
	Rules map[string]Filter
	// Inputs define the FilterInputs to parse
	// from
	Inputs []*model.FilterInput
}

type Filter struct {
	Join       string
	Conditions []model.FilterCondition
	Query      interface{}
}

// Parse FilterRules' Inputs to map onto a SQL Query
func (f FilterRules) Parse(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	for _, filter := range f.Inputs {
		filter.Key = util.ToSnake(filter.Key, false)

		key := filter.Key
		filterValue, hasRule := f.Rules[key]
		conditions := filterValue.Conditions
		// Make sure that Filter key, condition, and value(s) are
		// valid
		if !isValidKey(key, f.Keys) {
			return db, fmt.Errorf("invalid filter key provided. valid filter keys are: %s", strings.Join(f.Keys, ", "))
		}
		if !filter.HasValidCondition(conditions) {
			return db, errInvalidFilterCondition(key, conditions)
		}
		if err := checkFilterValue(filter); err != nil {
			return nil, err
		}
		// Check if a rule was provided for Filter Key
		if hasRule {
			// If Filter has a custom Query provided then
			// run it instead and continue to the next Filter
			fn, ok := filterValue.Query.(func(context.Context, *gorm.DB, *model.FilterInput) (*gorm.DB, error))
			if ok {
				customQueryDB, err := fn(ctx, db, filter)
				if err != nil {
					return nil, err
				}
				if customQueryDB.Error != nil {
					return nil, customQueryDB.Error
				}
				db = customQueryDB
				continue
			}
			// Check if filter requires a Preload to execute
			if len(filterValue.Join) > 0 {
				db = db.Joins(filterValue.Join)
			}
		}

		// Generate SQL Query
		query := key + filter.Condition.SQL()
		// Parse Condition and append query to db transaction
		switch filter.Condition {
		case model.FilterConditionBetween:
			db = db.Where(query, filter.Values[0], filter.Values[1])
		case model.FilterConditionIn, model.FilterConditionNotIn:
			db = db.Where(query, filter.Values)
		case model.FilterConditionIsNull, model.FilterConditionIsNotNull:
			db = db.Where(query)
		case model.FilterConditionLike, model.FilterConditionILike, model.FilterConditionNotLike:
			str := filter.Value.(string)
			db = db.Where(query, "%"+str+"%")
		default:
			db = db.Where(query, filter.Value)
		}
	}
	return db, nil
}

func isValidKey(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func checkFilterValue(filter *model.FilterInput) error {
	switch filter.Condition {
	case model.FilterConditionBetween:
		if len(filter.Values) != 2 {
			return fmt.Errorf("`%s` requires an array with exactly two elements inside the `values` field", filter.Condition.String())
		}
	case model.FilterConditionIn, model.FilterConditionNotIn:
		if len(filter.Values) < 1 {
			return fmt.Errorf("`%s` requires an array with exactly one element inside the `values` field", filter.Condition.String())
		}
	case model.FilterConditionLike, model.FilterConditionILike, model.FilterConditionNotLike:
		if filter.Value == nil {
			return fmt.Errorf("`%s` requires a valid item in the `value` field", filter.Condition.String())
		}
		_, ok := filter.Value.(string)
		if !ok {
			return fmt.Errorf("`%s` requires a string in the `value` field", filter.Condition.String())
		}
	default:
		if filter.Value == nil {
			return fmt.Errorf("`%s` requires a valid item in the `value` field", filter.Condition.String())
		}
	}
	return nil
}
