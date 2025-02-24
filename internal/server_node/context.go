package server_node

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
	"wordofwisdom/pkg/protocol"
)

type ServerContext struct {
	Ctx  context.Context
	Conn net.Conn

	clientTimeout       time.Duration
	maxMessageSizeBytes int
}

func NewServerContext(ctx context.Context, conn net.Conn, maxMessageSizeBytes int, clientTimeout time.Duration) *ServerContext {
	return &ServerContext{
		Ctx:                 ctx,
		Conn:                conn,
		maxMessageSizeBytes: maxMessageSizeBytes,
		clientTimeout:       clientTimeout,
	}
}

var (
	ErrConnectionClosed    = errors.New("connection closed")
	ErrFailedToReadMessage = errors.New("failed to read message")
	ErrFailedToSendMessage = errors.New("failed to send message")
	ErrClientTimeout       = errors.New("client timeout")
)

func (ctx *ServerContext) WaitMessage() (*protocol.RawMessage, error) {
	ctx.Conn.SetReadDeadline(time.Now().Add(ctx.clientTimeout))
	messageBuff := make([]byte, ctx.maxMessageSizeBytes)

	log.Printf("Waiting for message from client: %s for %s", ctx.Conn.RemoteAddr(), ctx.clientTimeout)

	bytesMessage, err := ctx.Conn.Read(messageBuff)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, ErrConnectionClosed
		}
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, ErrClientTimeout
		}
		return nil, errors.Join(err, ErrFailedToReadMessage)
	}
	log.Printf("Received message from client. [SIZE: %d bytes]", bytesMessage)

	return protocol.ParseRawMessage(messageBuff[:bytesMessage])
}

func (ctx *ServerContext) SendSuccessMessage(opcode uint32, msg protocol.MessageEncoder) error {
	return ctx.sendMessage(true, opcode, msg)
}

func (ctx *ServerContext) SendFailMessage(opcode uint32, msg protocol.MessageEncoder) error {
	return ctx.sendMessage(false, opcode, msg)
}

func (ctx *ServerContext) SendError(opcode uint32) error {
	return ctx.sendMessage(false, opcode, nil)
}

func (ctx *ServerContext) sendMessage(success bool, opcode uint32, payload protocol.MessageEncoder) error {
	rawMessage, err := protocol.BuildRawMessage(success, opcode, payload)
	if err != nil {
		return err
	}

	_, err = ctx.Conn.Write(rawMessage)
	if err != nil {
		return errors.Join(err, ErrFailedToSendMessage)
	}
	log.Printf("Sent message to client. [SIZE: %d bytes]", len(rawMessage))

	return nil
}
