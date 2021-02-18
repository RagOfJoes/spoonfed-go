package database

import (
	"context"
	"errors"

	"github.com/RagOfJoes/spoonfed-go/cmd/spoonfed-go/config"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Move the decode helper fn to FieldResolvers
// so that they just have a simple method that functions
// the same way.

var (
	recipeCollectionName = config.DatabaseCollectionNames["Recipe"]
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
func (db *DB) GetRecipes(ctx context.Context, limit int, cursor *string, sort *model.CursorSortInput) (*model.RecipeConnection, error) {
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
	filter := model.CursorToBson(cursor, key, order)
	cur, err := collection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$sort": decodedSort},
		{"$limit": int64(limit + 1)},
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
