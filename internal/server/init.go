package server

import (
	"github.com/RagOfJoes/spoonfed-go/internal/orm"
	"github.com/RagOfJoes/spoonfed-go/pkg/auth"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

// InitializeORM initiaalizes the ORM singleton object
func InitializeORM(cfg *util.ServerConfig) (*orm.ORM, error) {
	return orm.New(cfg)
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
