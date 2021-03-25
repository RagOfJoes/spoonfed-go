package resolver

import (
	"github.com/RagOfJoes/spoonfed-go/internal/orm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver base struct
type Resolver struct {
	ORM *orm.ORM
}
