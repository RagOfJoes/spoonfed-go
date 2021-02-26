package database

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/RagOfJoes/spoonfed-go/cmd/spoonfed-go/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Internal Properties
	connectOnce          sync.Once
	likeCollectionName   = config.DatabaseCollectionNames["Like"]
	userCollectionName   = config.DatabaseCollectionNames["User"]
	recipeCollectionName = config.DatabaseCollectionNames["Recipe"]

	// Singleton object
	db *DB

	// ErrClientNotInitialized defines an error where the mongo client has not yet
	// been initialized
	ErrClientNotInitialized = errors.New("mongo client has not been initialized")
	// ErrDatabaseNotSet defines an error where the database for DB
	// has not been set yet
	ErrDatabaseNotSet = errors.New("no database has been set")
	// ErrInvalidDatabaseName defines an error where a database name provided
	// does not exist within a Mongo Client
	ErrInvalidDatabaseName = errors.New("invalid database name provided")
	// ErrInvalidCollectionName defines an error where a collection name provided
	// does not exist within a Mongo Database
	ErrInvalidCollectionName = errors.New("invalid collection name provided")
)

// New instantiates a new db object that will be
func New(uri string) (*DB, error) {
	connectOnce.Do(func() {
		opts := options.Client()
		opts.ApplyURI(uri)
		opts.SetMaxPoolSize(5)
		client, err := mongo.NewClient(opts)
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		log.Print("[MongoDB] Successfully Connected")
		db = &DB{
			Client:      client,
			Collections: map[string]*mongo.Collection{},
		}
	})
	return db, nil
}

// Client allows the db singleton object to be accessed
func Client() (*DB, error) {
	if db.Client == nil {
		return nil, ErrClientNotInitialized
	}
	return db, nil
}

// SetDatabase sets database to singleton object
func (d *DB) SetDatabase(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	names, err := d.Client.ListDatabaseNames(ctx, primitive.D{})
	if err != nil {
		return err
	}
	contains := false
	for _, str := range names {
		if strings.EqualFold(str, name) {
			contains = true
		}
	}
	if !contains {
		return ErrInvalidDatabaseName
	}
	db.Database = d.Client.Database(name)
	return nil
}

// SetCollection sets collection to singleton object
func (d *DB) SetCollection(names ...string) error {
	if d.Database == nil {
		return ErrDatabaseNotSet
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collNames, err := d.Database.ListCollectionNames(ctx, primitive.D{})
	if err != nil {
		return err
	}
	for _, name := range collNames {
		if contains(names, name) {
			db.Collections[name] = db.Database.Collection(name)
		}
	}
	return nil
}

// Collection retrieves a collection with the given name
// from singleton object
func (d *DB) Collection(name string) (*mongo.Collection, error) {
	coll := d.Collections[name]
	if coll == nil {
		return nil, ErrInvalidCollectionName
	}
	return d.Collections[name], nil
}

// Simple helper func that will search to see if an array
// of strings contains a specific string
func contains(arr []string, str string) bool {
	for _, item := range arr {
		if strings.EqualFold(str, item) {
			return true
		}
	}
	return false
}
