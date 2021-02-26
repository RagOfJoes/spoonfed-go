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

// GetRecipes helper fn
func (db *DB) GetRecipes(ctx context.Context, limit int, cursor *string, sort *model.CursorSortInput, filters []*model.RecipeFilterInput) (*model.RecipeConnection, error) {
	if limit <= 1 || limit > 100 {
		return nil, errors.New("limit must be between the range of 1 - 100")
	}
	collection, err := db.Collection(recipeCollectionName)
	if err != nil {
		return nil, err
	}
	key := sort.Key()
	decodedSort := sort.Bson()
	order := sort.Value().Int()
	opts := options.Find()
	opts.SetLimit(int64(limit + 1))
	opts.SetSort(decodedSort)
	cursorFilter := model.CursorToBson(cursor, key, order)
	matchElements := bson.D{}
	if cursorFilter.Key != "" {
		matchElements = append(matchElements, cursorFilter)
	}
	for _, filter := range filters {
		if filter != nil {
			obj, err := filter.Bson()
			if err == nil {
				matchElements = append(matchElements, *obj)
			}
		}
	}
	cur, err := collection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchElements}},
		bson.D{{Key: "$sort", Value: decodedSort}},
		bson.D{{Key: "$limit", Value: int64(limit + 1)}},
	})
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
		if key == "date.creation" {
			unix := int(edge.Node.Date.Creation.Unix())
			edge.Cursor = model.EncodeCursor(unix)
		} else {
			edge.Cursor = model.EncodeCursor(edge.Node.Name)
		}
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
