package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/RagOfJoes/spoonfed-go/internal/auth"
	"github.com/RagOfJoes/spoonfed-go/internal/database"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/dataloader"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

func (r *queryResolver) GetRecipeDetail(ctx context.Context, slug string) (*model.Recipe, error) {
}

		return nil, err
	}
}

func (r *recipeResolver) NumOfLikes(ctx context.Context, obj *model.Recipe) (*int, error) {
	if err != nil {
		return nil, err
	}
}

func (r *recipeResolver) CreatedBy(ctx context.Context, obj *model.Recipe) (*model.User, error) {
	if err != nil {
	}
	return user, err
}

func (r *recipeResolver) IsLiked(ctx context.Context, obj *model.Recipe) (*bool, error) {
	if err != nil {
		return nil, err
	}
	}
}

// Recipe returns generated.RecipeResolver implementation.
func (r *Resolver) Recipe() generated.RecipeResolver { return &recipeResolver{r} }

type recipeResolver struct{ *Resolver }
