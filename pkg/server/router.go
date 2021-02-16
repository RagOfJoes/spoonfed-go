package server

import (
	"github.com/RagOfJoes/spoonfed-go/pkg/server/routes"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// AttachRoutes attaches routes to gin router
func AttachRoutes(cfg *util.ServerConfig, r *gin.Engine) error {
	if err := routes.GraphQL(cfg, r); err != nil {
		return err
	}
	return nil
}
