package database

import (
	"context"
	"log"
	"time"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	userCollectionName = "users"
)

// FindUserByID does exactly as the name suggests.
func (db *DB) FindUserByID(id primitive.ObjectID) (*model.User, error) {
	collection, err := db.Collection(userCollectionName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
func (db *DB) FindUsersByID(ids []string) ([]*model.User, []error) {
	collection, err := db.Collection(userCollectionName)
	if err != nil {
		return nil, []error{err}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
