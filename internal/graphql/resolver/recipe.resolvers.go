package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/dataloader"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"github.com/RagOfJoes/spoonfed-go/internal/orm/services"
	"github.com/gofrs/uuid"
)

func (r *mutationResolver) ToggleRecipeLike(ctx context.Context, recipeID string) (*model.Recipe, error) {
	return services.ToggleRecipeLike(ctx, r.ORM.DB, recipeID)
}

func (r *mutationResolver) EditRecipe(ctx context.Context, recipe model.EditRecipeInput) (*model.Recipe, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateRecipe(ctx context.Context, recipe model.NewRecipeInput) (*model.Recipe, error) {
	return services.CreateRecipe(ctx, r.ORM.DB, recipe)
}

func (r *queryResolver) GetRecipeDetail(ctx context.Context, slug string) (*model.Recipe, error) {
	return services.GetRecipeDetail(ctx, r.ORM.DB, slug)
}

func (r *queryResolver) GetRecipes(ctx context.Context, limit int, cursor *string, sort model.SortInput, filters []*model.FilterInput) (*model.RecipeConnection, error) {
	return services.GetRecipes(ctx, r.ORM.DB, limit, cursor, sort, filters)
}

func (r *recipeResolver) Images(ctx context.Context, obj *model.Recipe) ([]*model.Image, error) {
	tx := r.ORM.DB
	images := []models.Image{}
	if err := tx.Find(&images, "entity_type= ? AND entity_id = ?", "recipes", obj.ID).Error; err != nil {
		return nil, err
	}
	res := []*model.Image{}
	for _, image := range images {
		ok, err := image.IsValid(false)
		if err != nil || !ok {
			return nil, err
		}
		built, err := model.BuildImage(&image)
		if err != nil {
			return nil, err
		}
		res = append(res, built)
	}
	return res, nil
}

func (r *recipeResolver) NumOfLikes(ctx context.Context, obj *model.Recipe) (*int, error) {
	count, err := dataloader.For(ctx).RecipeLikeByID.Load(obj.ID)
	if err != nil {
		return nil, err
	}
	cast := int(*count)
	return &cast, nil
}

func (r *recipeResolver) CreatedBy(ctx context.Context, obj *model.Recipe) (*model.User, error) {
	user, err := dataloader.For(ctx).UserByID.Load(obj.CreatedBy)
	if err != nil {
		return services.GetUserFromDB(r.ORM.DB, obj.CreatedBy)
	}
	return user, err
}

func (r *recipeResolver) IsLiked(ctx context.Context, obj *model.Recipe) (*bool, error) {
	flag := false
	user, err := services.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &flag, nil
	}

	tx := r.ORM.DB.Begin()
	like := &models.Like{}
	tx.Find(like, &models.Like{EntityID: uuid.FromStringOrNil(obj.ID), EntityType: "recipes", UserID: user.ID})
	if err := tx.Error; err != nil {
		return nil, err
	}
	if like.UserID == user.ID && like.EntityID == uuid.FromStringOrNil(obj.ID) {
		if like.Active > 0 {
			flag = true
		}
	}
	return &flag, nil
}

// Recipe returns generated.RecipeResolver implementation.
func (r *Resolver) Recipe() generated.RecipeResolver { return &recipeResolver{r} }

type recipeResolver struct{ *Resolver }
