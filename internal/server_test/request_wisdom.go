package servertest

import (
	"context"
	"time"
	"wordofwisdom/internal/client_node"
	"wordofwisdom/internal/client_node/client_context"
	"wordofwisdom/internal/client_node/usecases"
	"wordofwisdom/pkg/server_sdk"
)

func RequestWisdomTest(ctx context.Context, cfg *client_node.ClientConfig) error {
	sdk := server_sdk.NewServerSDK(
		ctx,
		cfg.ServerAddress,
		cfg.MaxMessageSizeBytes,
		time.Duration(cfg.PopMessageTimeoutMs)*time.Millisecond,
	)
	if err := sdk.OpenConnection(); err != nil {
		return err
	}
	defer sdk.CloseConnection()

	clientCtx := client_context.NewClientContext(ctx, sdk)
	if err := usecases.RequestWisdom(clientCtx); err != nil {
		return err
	}

	return nil
}
