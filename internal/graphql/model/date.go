package model

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
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
