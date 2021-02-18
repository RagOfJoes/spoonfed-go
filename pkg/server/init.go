package server

import (
	"github.com/RagOfJoes/spoonfed-go/internal/auth"
	"github.com/RagOfJoes/spoonfed-go/internal/database"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

// InitializeDatabase initializes the db singleton object
func InitializeDatabase(cfg *util.ServerConfig) (*database.DB, error) {
	uri := cfg.Database.URI
	db, err := database.New(uri)
	if err != nil {
		return nil, err
	}
	if err := db.SetDatabase(cfg.Database.Name); err != nil {
		return nil, err
	}
	colls := make([]string, len(cfg.Database.Collections))
	for _, col := range cfg.Database.Collections {
		colls = append(colls, col)
	}
	if err := db.SetCollection(colls...); err != nil {
		return nil, err
	}
	return db, nil
}

// InitializeOpenIDClient does exactly as the name suggests
func InitializeOpenIDClient(cfg *util.ServerConfig) (*auth.OpenIDClient, error) {
	scope := cfg.Auth.Scopes
	clientID := cfg.Auth.ClientID
	clientSecret := cfg.Auth.ClientSecret
	openIDAutoDiscoveryURL := cfg.Auth.Issuer + "/.well-known/openid-configuration"
	client, err := auth.New(clientID, clientSecret, openIDAutoDiscoveryURL, scope)
	if err != nil {
		return nil, err
	}
	return client, nil
}
