package server

import (
	"context"
	"encoding/binary"
	"io"
	"wordofwisdom/internal/pow"
	"wordofwisdom/internal/protocol"
)

type serverHandlers struct {
	challengeDifficulty int
}

func NewServerHandlers(challengeDifficulty int) *serverHandlers {
	return &serverHandlers{
		challengeDifficulty: challengeDifficulty,
	}
}

func (h *serverHandlers) handleRequestWisdom(ctx context.Context, w io.Writer, r io.Reader, message []byte) error {
	challenge := pow.GenerateChallenge(h.challengeDifficulty)
	quote := GetRandomQuote()
	messageBuff := make([]byte, 4)
	binary.BigEndian.PutUint32(messageBuff, protocol.OPCODE_REQUEST_WISDOM)
	w.Write(messageBuff)
	w.Write([]byte(quote))
	return nil
}
