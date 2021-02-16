package routes

import (
	"log"

	"github.com/RagOfJoes/spoonfed-go/pkg/server/handlers"
	"github.com/RagOfJoes/spoonfed-go/pkg/server/middlewares"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// GraphQL sets up GraphQL related routes
func GraphQL(cfg *util.ServerConfig, r *gin.Engine) error {
	// GraphQL paths
	graphqlPath := cfg.GraphQL.Path
	playgroundPath := cfg.GraphQL.PlaygroundPath

	// GraphQL handler
	// Middlewares execution order:
	// 1. Auth
	// 2. Dataloaders
	r.POST(graphqlPath, middlewares.Auth(graphqlPath), middlewares.Dataloader(), handlers.GraphQLHandler(&cfg.GraphQL))
	log.Printf("[GraphQL] mounted at %s", graphqlPath)
	// Playground handler
	if cfg.GraphQL.EnablePlayground {
		log.Printf("[GraphQL-Playground] mounted at %s", playgroundPath)
		r.GET(playgroundPath, handlers.PlaygroundHandler(playgroundPath))
	}
	return nil
}
