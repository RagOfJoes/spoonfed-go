package main

import (
	"log"

	"github.com/RagOfJoes/spoonfed-go/cmd/spoonfed-go/config"
	"github.com/RagOfJoes/spoonfed-go/pkg/server"
)

func main() {
	cfg := config.LoadConfig()
	err := server.Run(cfg)
	if err != nil {
		log.Panic(err)
	}
}
