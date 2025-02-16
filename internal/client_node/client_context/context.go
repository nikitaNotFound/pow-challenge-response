package client_context

import (
	"context"
	"wordofwisdom/pkg/server_sdk"
)

type ClientContext struct {
	Ctx context.Context
	Sdk *server_sdk.ServerSDK
}

func NewClientContext(ctx context.Context, sdk *server_sdk.ServerSDK) *ClientContext {
	return &ClientContext{
		Ctx: ctx,
		Sdk: sdk,
	}
}
