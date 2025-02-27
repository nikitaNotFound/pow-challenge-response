package servertest

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"wordofwisdom/internal/client_node"
)

func RunTests(ctx context.Context) {
	cfg := client_node.GetClientConfig()

	successCounter := atomic.Int32{}
	errCounter := atomic.Int32{}
	errContainer := make([]error, 0, 30)

	wg := sync.WaitGroup{}
	for range 100 {
		wg.Add(1)
		go func() error {
			defer wg.Done()
			if err := RequestWisdomTest(ctx, cfg); err != nil {
				errCounter.Add(1)
				errContainer = append(errContainer, err)
				return err
			}

			successCounter.Add(1)

			return nil
		}()
	}

	wg.Wait()

	log.Printf("Successfully completed %d tests. [errors: %d]", successCounter.Load(), errCounter.Load())

	if errCounter.Load() > 0 {
		log.Printf("Errors: %v", errContainer)
	}
}
