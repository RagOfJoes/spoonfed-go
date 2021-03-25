package model

import (
	"fmt"
	"io"
	"strconv"
)

type FilterCondition string

const (
	FilterConditionEquals           FilterCondition = "EQUALS"
	FilterConditionNotEquals        FilterCondition = "NOT_EQUALS"
	FilterConditionLessThan         FilterCondition = "LESS_THAN"
	FilterConditionGreaterThan      FilterCondition = "GREATER_THAN"
	FilterConditionLessThanEqual    FilterCondition = "LESS_THAN_EQUAL"
	FilterConditionGreaterThanEqual FilterCondition = "GREATER_THAN_EQUAL"
	FilterConditionBetween          FilterCondition = "BETWEEN"
	FilterConditionIs               FilterCondition = "IS"
	FilterConditionIsNull           FilterCondition = "IS_NULL"
	FilterConditionIsNotNull        FilterCondition = "IS_NOT_NULL"
	FilterConditionIn               FilterCondition = "IN"
	FilterConditionNotIn            FilterCondition = "NOT_IN"
	FilterConditionLike             FilterCondition = "LIKE"
	FilterConditionILike            FilterCondition = "ILIKE"
	FilterConditionNotLike          FilterCondition = "NOT_LIKE"
)

var (
	filterConditionMap = map[FilterCondition]string{
		FilterConditionEquals:           " = ?",
		FilterConditionNotEquals:        " != ?",
		FilterConditionLessThan:         " < ?",
		FilterConditionGreaterThan:      " > ?",
		FilterConditionLessThanEqual:    " <= ?",
		FilterConditionGreaterThanEqual: " >= ?",
		FilterConditionBetween:          " BETWEEN ? AND ?",
		FilterConditionIs:               " IS ?",
		FilterConditionIsNull:           " IS NULL ?",
		FilterConditionIsNotNull:        " IS NOT NULL ?",
		FilterConditionIn:               " IN (?)",
		FilterConditionNotIn:            " NOT IN (?)",
		FilterConditionLike:             " LIKE ?",
		FilterConditionILike:            " ILIKE ?",
		FilterConditionNotLike:          " NOT LIKE ?",
	}
)

var AllFilterCondition = []FilterCondition{
	FilterConditionEquals,
	FilterConditionNotEquals,
	FilterConditionLessThan,
	FilterConditionGreaterThan,
	FilterConditionLessThanEqual,
	FilterConditionGreaterThanEqual,
	FilterConditionBetween,
	FilterConditionIs,
	FilterConditionIsNull,
	FilterConditionIsNotNull,
	FilterConditionIn,
	FilterConditionNotIn,
	FilterConditionLike,
	FilterConditionILike,
	FilterConditionNotLike,
}

func (e FilterCondition) IsValid() bool {
	switch e {
	case FilterConditionEquals,
		FilterConditionNotEquals,
		FilterConditionLessThan,
		FilterConditionGreaterThan,
		FilterConditionLessThanEqual,
		FilterConditionGreaterThanEqual,
		FilterConditionBetween,
		FilterConditionIs,
		FilterConditionIsNull,
		FilterConditionIsNotNull,
		FilterConditionIn,
		FilterConditionNotIn,
		FilterConditionLike,
		FilterConditionILike,
		FilterConditionNotLike:
		return true
	}
	return false
}

func (e FilterCondition) String() string {
	return string(e)
}

func (e FilterCondition) SQL() string {
	return filterConditionMap[e]
}

func (e *FilterCondition) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterCondition", str)
	}
	return nil
}

func (e FilterCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
