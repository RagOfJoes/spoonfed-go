package database

import "go.mongodb.org/mongo-driver/mongo"

// Collections Type
type Collections struct {
	Recipes *mongo.Collection
	Users   *mongo.Collection
}

// DB Type
type DB struct {
	Client      *mongo.Client
	Database    *mongo.Database
	Databases   map[string]*mongo.Database
	Collections map[string]*mongo.Collection
}
