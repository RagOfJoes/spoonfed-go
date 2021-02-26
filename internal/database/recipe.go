package database

import (
	"context"
	"errors"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"go.mongodb.org/mongo-driver/bson"
)

// GetRecipeDetail helper fn
func (db *DB) GetRecipeDetail(ctx context.Context, slug string) (*model.Recipe, error) {
	if slug == "" {
		return nil, errors.New("must provide valid slug")
	}
	collection, err := db.Collection(recipeCollectionName)
	if err != nil {
		return nil, err
	}
	var recipe *model.Recipe
	fErr := collection.FindOne(ctx, bson.D{{Key: "slug", Value: slug}}).Decode(&recipe)
	if fErr != nil {
		return nil, err
	}
	return recipe, nil
}
