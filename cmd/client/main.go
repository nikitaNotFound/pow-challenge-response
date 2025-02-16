package main

import (
	"context"
	"log"
	"wordofwisdom/internal/client_node"
)

func main() {
	ctx := context.Background()

	if err := client_node.RunClient(ctx); err != nil {
		log.Fatalf("Failed to run client: %v", err)
	}
}
