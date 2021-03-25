package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidMetaDate = errors.New("invalid time type provided")
)

// MetaDate is the custom type for MetaDate scalar
type MetaDate struct {
	Creation   time.Time
	LastUpdate *time.Time
}

func (m *MetaDate) IsValid() (ok bool, err error) {
	if m.Creation.IsZero() {
		return false, ErrInvalidMetaDate
	}
	return true, nil
}
