package server_sdk

import (
	"context"
	"errors"
	"log"
	"net"
	"wordofwisdom/pkg/protocol"
)

type ServerSDK struct {
	serverAddress string
	ctx           context.Context

	conn net.Conn

	maxMessageSizeBytes int
}

func NewServerSDK(ctx context.Context, address string, maxMessageSizeBytes int) *ServerSDK {
	return &ServerSDK{
		serverAddress:       address,
		ctx:                 ctx,
		maxMessageSizeBytes: maxMessageSizeBytes,
	}
}

var (
	ErrConnectionClosed     = errors.New("connection closed")
	ErrConnectionFailed     = errors.New("connection failed")
	ErrMessageTooShort      = errors.New("message is too short")
	ErrFailedToWaitMessage  = errors.New("failed to wait message")
	ErrFailedToSendMessage  = errors.New("failed to send message")
	ErrFailedToBuildMessage = errors.New("failed to build message")
)

func (s *ServerSDK) OpenConnection() error {
	conn, err := net.Dial("tcp", s.serverAddress)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return ErrConnectionClosed
		}
		return errors.Join(err, ErrConnectionFailed)
	}
	s.conn = conn

	return nil
}

func (s *ServerSDK) CloseConnection() error {
	return s.conn.Close()
}

func (s *ServerSDK) SendMessage(success bool, opcode uint32, payload protocol.MessageEncoder) error {
	rawMessage, err := protocol.BuildRawMessage(success, opcode, payload)
	if err != nil {
		return errors.Join(err, ErrFailedToBuildMessage)
	}

	_, err = s.conn.Write(rawMessage)
	if err != nil {
		return errors.Join(err, ErrFailedToSendMessage)
	}

	return nil
}

func (s *ServerSDK) WaitMessage() (*protocol.RawMessage, error) {
	messageBuff := make([]byte, s.maxMessageSizeBytes)
	bytesMessage, err := s.conn.Read(messageBuff)
	log.Printf("Received message from server, %d bytes", bytesMessage)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return nil, ErrConnectionClosed
		}
		return nil, errors.Join(err, ErrFailedToWaitMessage)
	}

	return protocol.ParseRawMessage(messageBuff[:bytesMessage])
}
