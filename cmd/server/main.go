package main

import (
	"context"
	"log"
	"wordofwisdom/internal/server_node"
)

func main() {
	ctx := context.Background()
	if err := server_node.RunServer(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
