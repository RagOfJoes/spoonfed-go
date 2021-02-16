package database

import "go.mongodb.org/mongo-driver/mongo"

// DB Type
type DB struct {
	Client      *mongo.Client
	Database    *mongo.Database
	Collections map[string]*mongo.Collection
}
