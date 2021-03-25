package routes

import (
	"github.com/RagOfJoes/spoonfed-go/internal/orm"
	"github.com/RagOfJoes/spoonfed-go/internal/server/handlers"
	"github.com/RagOfJoes/spoonfed-go/internal/server/middlewares"
	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// GraphQL sets up GraphQL related routes
func GraphQL(cfg *util.ServerConfig, r *gin.Engine, o *orm.ORM) error {
	// GraphQL paths
	graphqlPath := cfg.GraphQL.Path
	playgroundPath := cfg.GraphQL.PlaygroundPath

	// GraphQL handler
	// Middlewares execution order:
	// 1. Auth
	// 2. Dataloaders
	r.POST(graphqlPath, middlewares.Auth(graphqlPath, o), middlewares.Dataloader(o), handlers.GraphQLHandler(&cfg.GraphQL, o))
	logger.Infof("[GraphQL] mounted at %s", graphqlPath)
	// Playground handler
	if cfg.GraphQL.EnablePlayground {
		logger.Infof("[GraphQL-Playground] mounted at %s", playgroundPath)
		r.GET(playgroundPath, handlers.PlaygroundHandler(playgroundPath))
	}
	return nil
}
