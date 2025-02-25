package main

import (
	"context"
	tests "wordofwisdom/internal/server_test"
)

func main() {
	ctx := context.Background()

	tests.RunTests(ctx)
}
