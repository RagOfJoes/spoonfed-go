package database

import (
	"context"
	"errors"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"go.mongodb.org/mongo-driver/bson"
)

// GetRecipes helper fn
func (db *DB) GetRecipes(ctx context.Context, limit int, cursor *string, sort *model.RecipeSortInput, filters []*model.RecipeFilterInput) (*model.RecipeConnection, error) {
	if limit <= 1 || limit > 100 {
		return nil, errors.New("limit must be between the range of 1 - 100")
	}
	collection, err := db.Collection(recipeCollectionName)
	if err != nil {
		return nil, err
	}
	sortOrder := sort.Value().Int()
	cursorBson := model.CursorToBson(cursor, sort.BsonKey(), sortOrder)
	inputs := []model.FilterInput{}
	factory := model.RecipeFilterInputFactory{}
	for _, f := range filters {
		inputs = append(inputs, factory.Create(*f))
	}
	filterPipe := model.FilterPipeline{
		Filters: inputs,
		Prerequisites: map[string][]bson.M{
			"user.username": {
				{
					"$lookup": bson.M{
						"foreignField": "sub",
						"as":           "user",
						"localField":   "createdBy",
						"from":         userCollectionName,
					},
				},
			},
		},
	}
	sortPipe := sort.Bson()
	match := filterPipe.BuildMatch(cursorBson)
	lookup := filterPipe.BuildPreReq(sortPipe, match)
	pipeline := []bson.M{}
	if len(lookup) > 0 {
		for _, look := range lookup {
			pipeline = append(pipeline, look)
		}
	}
	pipeline = append(pipeline,
		bson.M{"$match": match},
		bson.M{"$sort": sortPipe},
		bson.M{"$limit": int64(limit + 1)},
	)
	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	// Iterate through query result
	var edges []*model.RecipeEdge
	for cur.Next(ctx) {
		// Initialize RecipeEdge
		edge := &model.RecipeEdge{
			Cursor: "",
			Node:   &model.Recipe{},
		}
		// Unmarshal document to Recipe
		err := cur.Decode(edge.Node)
		if err != nil {
			return nil, err
		}
		// Assign Edge a cursor
		edge.Cursor = sort.Cursor(edge.Node)
		// Append Edge to Edges
		edges = append(edges, edge)
	}
	// Check for errors
	if err := cur.Err(); err != nil {
		return nil, err
	}
	// Handle new PageInfo
	var newCursor string
	hasNextPage := len(edges) > limit-1
	if hasNextPage {
		newCursor = edges[len(edges)-1].Cursor
		edges = edges[:len(edges)-1]
	} else {
		newCursor = ""
	}
	newPageInfo := &model.PageInfo{
		Cursor:      newCursor,
		HasNextPage: hasNextPage,
	}
	// Return new RecipeConnection
	connection := &model.RecipeConnection{
		Edges:    edges,
		PageInfo: newPageInfo,
	}
	return connection, nil
}
