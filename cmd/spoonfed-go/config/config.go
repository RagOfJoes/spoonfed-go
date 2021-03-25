package config

import (
	"strings"

	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/joho/godotenv"
)

// LoadConfig does just that
func LoadConfig() *util.ServerConfig {
	err := godotenv.Load()
	if err != nil && util.GetAssert("APP_ENV") != "PRODUCTION" {
		logger.Panic("[ENV] Failed to load")
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
		ORM: util.ORMConfig{
			DSN:         util.GetAssert("ORM_DSN"),
			AutoMigrate: util.GetAssertBool("ORM_AUTOMIGRATE"),
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
