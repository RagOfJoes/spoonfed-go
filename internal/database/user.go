package database

import (
	"context"
	"log"

	"github.com/RagOfJoes/spoonfed-go/cmd/spoonfed-go/config"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userCollectionName = config.DatabaseCollectionNames["User"]
)

// FindUserByID does exactly as the name suggests.
func (db *DB) FindUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	collection, err := db.Collection(userCollectionName)
	if err != nil {
		return nil, err
	}
	var user *model.User
	fErr := collection.FindOne(ctx, bson.D{{Key: "sub", Value: id}}).Decode(&user)
	if fErr != nil {
		log.Print(err)
	}

	return user, nil
}

// FindUsersByID does exactly as the name suggests.
// This is mainly used for the Dataloader's batch
// function
func (db *DB) FindUsersByID(ctx context.Context, ids []string) ([]*model.User, []error) {
	collection, err := db.Collection(userCollectionName)
	if err != nil {
		return nil, []error{err}
	}
	var users []*model.User
	var errors []error
	var converted []primitive.ObjectID
	for i := 0; i < len(ids); i++ {
		id, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			errors = append(errors, err)
			continue
		}
		converted = append(converted, id)
	}
	cur, err := collection.Find(ctx, bson.D{
		{
			Key: "sub", Value: bson.D{
				{
					Key: "$in", Value: converted,
				},
			},
		},
	})
	if err != nil {
		errors = append(errors, err)
	}
	for cur.Next(ctx) {
		user := &model.User{}
		err := cur.Decode(user)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		users = append(users, user)
	}
	return users, errors
}
