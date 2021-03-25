package dataloader

import (
	"context"

	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

// Loaders defines DataLoaders
type Loaders struct {
	UserByID       *UserLoader
	RecipeLikeByID *RecipeLikeLoader
}

// For loads existing Dataloaders provided a context
func For(ctx context.Context) *Loaders {
	return ctx.Value(util.ProjectContextKeys.Dataloader).(*Loaders)
}
