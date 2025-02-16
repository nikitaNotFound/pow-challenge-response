package server_node

import (
	"context"
	"wordofwisdom/pkg/protocol/requests"
)

func RunServer(ctx context.Context) error {
	cfg := GetServerConfig()

	tcpServer := NewTcpServer(ctx, cfg.Address, cfg.MaxMessageSizeBytes, cfg.MaxConnectionsPerClient)

	handlers := NewServerHandlers(cfg.ChallengeDifficulty)
	tcpServer.RegisterHandler(requests.OPCODE_REQUEST_WISDOM, handlers.handleRequestWisdom)

	return tcpServer.Run()
}
