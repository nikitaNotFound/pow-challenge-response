package client_node

import (
	"context"
	"fmt"
	"wordofwisdom/internal/client_node/client_context"
	"wordofwisdom/internal/client_node/usecases"
	"wordofwisdom/pkg/server_sdk"
)

func RunClient(ctx context.Context) error {
	cfg := GetClientConfig()

	sdk := server_sdk.NewServerSDK(ctx, cfg.ServerAddress, cfg.MaxMessageSizeBytes)
	if err := sdk.OpenConnection(); err != nil {
		return err
	}
	defer sdk.CloseConnection()

	clientCtx := client_context.NewClientContext(ctx, sdk)

	for {
		userInput := ""
		fmt.Scanln(&userInput)

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
