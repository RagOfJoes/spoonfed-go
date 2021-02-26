package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecipeLikes returns the number of likes that
// a recipe has
func (db *DB) RecipeLikes(ctx context.Context, id string) int {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0
	}
	collection, err := db.Collection(likeCollectionName)
	if err != nil {
		return 0
	}
	counts, err := collection.CountDocuments(ctx, bson.D{
		{Key: "active", Value: true},
		{Key: "recipeId", Value: objID},
	})
	if err != nil {
		return 0
	}
	count := int(counts)
	return count
}

// IsRecipeLiked returns whether a recipe is liked by a user
func (db *DB) IsRecipeLiked(ctx context.Context, userID string, recipeID string) bool {
	user, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false
	}
	recipe, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		return false
	}
	collection, err := db.Collection(likeCollectionName)
	if err != nil {
		return false
	}
	liked, err := collection.CountDocuments(ctx, bson.D{
		{Key: "active", Value: true},
		{Key: "userId", Value: user},
		{Key: "recipeId", Value: recipe},
	})
	if liked == 1 {
		return true
	}
	return false
}
