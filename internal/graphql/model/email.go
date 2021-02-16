package model

import (
	"errors"
	"regexp"

	"github.com/99designs/gqlgen/graphql"
)

// Email is the custom type for Email scalar
type Email string

var (
	// ErrEmailType error
	ErrEmailType = errors.New("email must be a string")
	// ErrEmailFormat error
	ErrEmailFormat = errors.New("invalid format")
	emailRegexp    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func (e Email) String() string {
	return string(e)
}

// IsValid fn checks whether a string matches our custom
// email regex
func (e Email) IsValid() bool {
	if !emailRegexp.MatchString(e.String()) {
		return false
	}
	return true
}

// MarshalEmail implements graphql.Marshaller into scalar
func MarshalEmail(e Email) graphql.Marshaler {
	if e.IsValid() {
		return graphql.MarshalString(e.String())
	}
	return graphql.Null
}

// UnmarshalEmail implements graphql.Unmarshaller into scalar
func UnmarshalEmail(v interface{}) (Email, error) {
	str, ok := v.(string)
	if !ok {
		return "null", ErrEmailType
	}

	e := Email(str)
	if !e.IsValid() {
		return "null", ErrEmailFormat
	}
	return e, nil
}
