package model

import (
	"fmt"
	"io"
	"strconv"
)

// SortOrder defines the SortOrder type that correlates
// to the graphql schamas.
type SortOrder string

const (
	// SortOrderAsc defines the ascending order.
	SortOrderAsc SortOrder = "ASC"
	// SortOrderDesc defines the descending order.
	SortOrderDesc SortOrder = "DESC"
)

// IsValid ensures that the value provided
// is within the enum's definition.
func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

// String converts the value into a string.
func (e SortOrder) String() string {
	return string(e)
}

// Int converts SortOrder value into an integer.
func (e SortOrder) Int() int {
	switch e {
	case SortOrderAsc:
		return 1
	default:
		return -1
	}
}

// UnmarshalGQL is a custom unmarshall fn for SortOrder type.
func (e *SortOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortOrder", str)
	}
	return nil
}

// MarshalGQL is a custom marshall fn for SortOrder type.
func (e SortOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
