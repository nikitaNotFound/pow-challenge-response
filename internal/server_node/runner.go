package server_node

import (
	"context"
	"net/http"
	"wordofwisdom/pkg/protocol/requests"
	_ "wordofwisdom/pkg/wrapper_expvars"
)

func RunServer(ctx context.Context) error {
	cfg := GetServerConfig()

	tcpServer := NewTcpServer(ctx, cfg.Address, cfg.MaxMessageSizeBytes, cfg.MaxConnectionsPerClient, cfg.WorkersAmount)

	handlers := NewServerHandlers(cfg.ChallengeDifficulty)
	tcpServer.RegisterHandler(requests.OPCODE_REQUEST_WISDOM, handlers.handleRequestWisdom)

	go http.ListenAndServe(":1234", nil)

	return tcpServer.Run()
}
