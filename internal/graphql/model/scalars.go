package model

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

// Date is the custom type for Date scalar
type Date time.Time

var (
	// ErrDateFormat error
	ErrDateFormat = errors.New("date must be RFC3339 formatted string")
)

// Time converts custome Date type to go's Time type
func (d Date) Time() time.Time {
	return time.Time(d)
}

// MarshalDate implements graphql.Marshaller into scalar
func MarshalDate(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(t.Format(time.RFC3339)))
	})
}

// UnmarshalDate implements graphql.Unmarshaller into scalar
func UnmarshalDate(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		return time.Parse(time.RFC3339, tmpStr)
	}
	return time.Time{}, ErrDateFormat
}

// Email is the custom type for Email scalar
type Email string

var (
	// ErrEmailType error
	ErrEmailType = errors.New("email must be a string")
	// ErrEmailFormat error
	ErrEmailFormat = errors.New("invalid format")
)

func (e Email) String() string {
	return string(e)
}

// IsValid fn checks whether a string matches our custom
// email regex
func (e Email) IsValid() bool {
	if !util.IsEmail(e.String()) {
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
