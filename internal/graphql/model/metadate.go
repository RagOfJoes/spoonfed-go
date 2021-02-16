package model

import (
	"time"
)

// MetaDate is the custom type for MetaDate scalar
type MetaDate struct {
	Creation   time.Time
	LastUpdate time.Time
}
