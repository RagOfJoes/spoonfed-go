package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/resolver"
	"github.com/RagOfJoes/spoonfed-go/internal/orm"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// GraphQLHandler configures gqlgen and returns a gin.HanlderFunc
// that can be attached to a route
func GraphQLHandler(cfg *util.GraphQLConfig, o *orm.ORM) gin.HandlerFunc {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		ORM: o,
	}})
	h := handler.New(schema)

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})
	h.SetQueryCache(lru.New(1000))

	// Extensions
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	if cfg.EnableIntrospection {
		h.Use(extension.Introspection{})
	}

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler configures a GraphQL plaground
// and returns a gin.HanlderFunc that can be attached
// to a route
func PlaygroundHandler(path string) gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", path)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
