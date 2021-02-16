package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe with ObjectID
type Recipe struct {
	ID           string             `json:"_id" bson:"_id"`
	Name         string             `json:"name"`
	Servings     string             `json:"servings"`
	Time         *RecipeTimeType    `json:"time"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	Slug         string             `json:"slug"`
	ImportURL    *string            `json:"importUrl"`
	Images       []*Image           `json:"images"`
	Date         *MetaDate          `json:"date"`
	CreatedBy    primitive.ObjectID `json:"user"`
}
