package main

import (
	"context"
	"log"
	"wordofwisdom/internal/server"
)

func main() {
	ctx := context.Background()
	if err := server.RunServer(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
