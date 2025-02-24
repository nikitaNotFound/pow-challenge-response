package client_node

import (
	"context"
	"fmt"
	"log"
	"wordofwisdom/internal/client_node/client_context"
	"wordofwisdom/internal/client_node/usecases"
	"wordofwisdom/pkg/server_sdk"
)

func RunClient(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cfg := GetClientConfig()

	sdk := server_sdk.NewServerSDK(ctx, cfg.ServerAddress, cfg.MaxMessageSizeBytes)
	if err := sdk.OpenConnection(); err != nil {
		return err
	}
	defer sdk.CloseConnection()

	clientCtx := client_context.NewClientContext(ctx, sdk)

	userInputCh := make(chan string)
	go func() {
		for {
			var userInput string
			fmt.Scanln(&userInput)
			userInputCh <- userInput
		}
	}()

	go func() {
		if err := sdk.WaitForClose(); err != nil {
			log.Printf("Connection closed by server: %v", err)
			cancel()
		}
	}()

	for {
		userInput := ""

		select {
		case <-ctx.Done():
			return nil
		case userInput = <-userInputCh:
		}

		if userInput == "exit" {
			break
		}

		if userInput == "wisdom" {
			if err := usecases.RequestWisdom(clientCtx); err != nil {
				return err
			}
		}
	}

	return nil
}
