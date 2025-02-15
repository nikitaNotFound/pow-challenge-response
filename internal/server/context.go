package server

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"wordofwisdom/internal/protocol"
)

type ServerContext struct {
	Ctx        context.Context
	Conn       net.Conn
	RawMessage []byte

	maxMessageSizeBytes int
}

func NewServerContext(ctx context.Context, conn net.Conn, maxMessageSizeBytes int) *ServerContext {
	return &ServerContext{
		Ctx:                 ctx,
		Conn:                conn,
		maxMessageSizeBytes: maxMessageSizeBytes,
	}
}

func (ctx *ServerContext) WaitMessage() (*protocol.Message, error) {
	messageBuff := make([]byte, ctx.maxMessageSizeBytes)
	bytesMessage, err := ctx.Conn.Read(messageBuff)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("connection closed by client")
		}
		return nil, errors.Join(err, errors.New("failed to read message"))
	}

	flags := messageBuff[0]
	opcode := binary.BigEndian.Uint32(messageBuff[1:5])
	return &protocol.Message{
		Flags:  flags,
		Opcode: opcode,
		Data:   messageBuff[5:bytesMessage],
	}, nil
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
	messageBuff := make([]byte, 5)

	flags := protocol.EmptyMessageFlags()
	if !success {
		flags.SetFlag(protocol.MSG_FAIL_FLAG)
	}
	messageBuff[0] = byte(flags)

	binary.BigEndian.PutUint32(messageBuff[1:5], opcode)

	if payload != nil {
		buff, err := payload.Encode()
		if err != nil {
			return errors.Join(err, errors.New("failed to encode message"))
		}

		messageBuff = append(messageBuff, buff...)
	}

	_, err := ctx.Conn.Write(messageBuff)
	if err != nil {
		return errors.Join(err, errors.New("failed to send message"))
	}

	return nil
}
