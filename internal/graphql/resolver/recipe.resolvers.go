package graphql

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

func (r *queryResolver) GetRecipes(ctx context.Context, limit int, cursor *string, sort *model.CursorSortInput) (*model.RecipeConnection, error) {
	client, err := database.Client()
	if err != nil {
		return nil, err
	}
	recipes, err := client.GetRecipes(ctx, limit, cursor, sort)
	return recipes, err
}

func (r *recipeResolver) NumOfLikes(ctx context.Context, obj *model.Recipe) (*int, error) {
	client, err := database.Client()
	if err != nil {
		return nil, err
	}
	count := client.RecipeLikes(ctx, obj.ID)
	return &count, nil
}

func (r *recipeResolver) CreatedBy(ctx context.Context, obj *model.Recipe) (*model.User, error) {
	user, err := dataloader.For(ctx).UserByID.Load(obj.CreatedBy.Hex())
	// If DataLoader fails then fallback to just using database
	if err != nil {
		client, err := database.Client()
		if err != nil {
			return nil, err
		}
		return client.FindUserByID(ctx, obj.CreatedBy)
	}
	return user, err
}

func (r *recipeResolver) IsLiked(ctx context.Context, obj *model.Recipe) (*bool, error) {
	client, err := database.Client()
	if err != nil {
		return nil, err
	}
	user, err := auth.GetUserFromContext(ctx, util.ProjectContextKeys.User)
	if err != nil || user == nil {
		return nil, nil
	}
	liked := client.IsRecipeLiked(ctx, user.Sub, obj.ID)
	return &liked, nil
}

// Recipe returns generated.RecipeResolver implementation.
func (r *Resolver) Recipe() generated.RecipeResolver { return &recipeResolver{r} }

type recipeResolver struct{ *Resolver }
