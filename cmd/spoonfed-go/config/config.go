package config

import (
	"log"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/joho/godotenv"
)

var (
	// DatabaseCollectionNames is just that
	DatabaseCollectionNames = map[string]string{
		"User":     "users",
		"Like":     "likes",
		"Recipe":   "recipes",
		"Creation": "creations",
	}
)

// LoadConfig does just that
func LoadConfig() *util.ServerConfig {
	err := godotenv.Load()
	if err != nil {
		log.Panic("[ENV] Failed to load")
		return nil
	}
	return &util.ServerConfig{
		Port:   util.GetAssert("PORT"),
		Host:   util.GetAssert("SERVER_HOST"),
		Scheme: util.GetAssert("SERVER_SCHEME"),
		GraphQL: util.GraphQLConfig{
			Path:                util.GetAssert("GRAPHQL_PATH"),
			PlaygroundPath:      util.GetAssert("GRAPHQL_PLAYGROUND_PATH"),
			EnablePlayground:    util.GetAssertBool("GRAPHQL_PLAYGROUND_ENABLE"),
			EnableIntrospection: util.GetAssertBool("GRAPHQL_INTROSPECTION_ENABLE"),
		},
		Database: util.DatabaseConfig{
			Collections: DatabaseCollectionNames,
			URI:         util.GetAssert("MONGO_URI"),
			Name:        util.GetAssert("MONGO_DB_NAME"),
		},
		Auth: util.OpenIDClientConfig{
			Issuer:       util.GetAssert("ROJ_ISSUER"),
			ClientID:     util.GetAssert("ROJ_CLIENT_ID"),
			ClientSecret: util.GetAssert("ROJ_CLIENT_SECRET"),
			// Should be structured as such: openid,profile,etc.
			Scopes: strings.Split(util.GetAssert("ROJ_SCOPES"), ","),
		},
	}
}
