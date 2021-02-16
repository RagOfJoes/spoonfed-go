package server

import (
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// Run initializes required modules and then runs the server
func Run(cfg *util.ServerConfig) error {
	if _, err := InitializeDatabase(cfg); err != nil {
		return err
	}
	if _, err := InitializeOpenIDClient(cfg); err != nil {
		return err
	}
	router := gin.Default()
	if err := AttachRoutes(cfg, router); err != nil {
		return err
	}
	if err := router.Run(cfg.ListenEndpoint()); err != nil {
		return err
	}
	return nil
}
