package server

import (
	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// Run initializes required modules and then runs the server
func Run(cfg *util.ServerConfig) error {
	orm, err := InitializeORM(cfg)
	if err != nil {
		logger.Error("[DB] Failed to initialize database.")
		return err
	}
	if _, err := InitializeOpenIDClient(cfg); err != nil {
		logger.Error("[Auth] Failed to initialize OIDC client.")
		return err
	}
	router := gin.Default()
	if err := AttachRoutes(cfg, router, orm); err != nil {
		logger.Error("[Server] Failed to attach routes.")
		return err
	}
	if err := router.Run(cfg.ListenEndpoint()); err != nil {
		logger.Error("[Server] Failed to run server.")
		return err
	}
	return nil
}
