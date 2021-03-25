package main

import (
	"github.com/RagOfJoes/spoonfed-go/cmd/spoonfed-go/config"
	"github.com/RagOfJoes/spoonfed-go/internal/server"
	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	err := server.Run(cfg)
	if err != nil {
		logger.Panic(err)
	}
}
