package main

import (
	"log"

	"naiimage/backend/internal/config"
	"naiimage/backend/internal/server"
)

func main() {
	cfg := config.Load()
	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
