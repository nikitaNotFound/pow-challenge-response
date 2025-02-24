package server_sdk

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"wordofwisdom/pkg/protocol"
)

type ServerSDK struct {
	serverAddress       string
	maxMessageSizeBytes int
	ctx                 context.Context

	conn net.Conn

	messagesCh  chan []byte
	connCloseCh chan error
	errCh       chan error
}

func NewServerSDK(ctx context.Context, address string, maxMessageSizeBytes int) *ServerSDK {
	return &ServerSDK{
		serverAddress:       address,
		ctx:                 ctx,
		maxMessageSizeBytes: maxMessageSizeBytes,
		messagesCh:          make(chan []byte),
		connCloseCh:         make(chan error),
		errCh:               make(chan error),
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

	go s.startReceivingMessages()

	return nil
}

func (s *ServerSDK) startReceivingMessages() {
	messageBuff := make([]byte, s.maxMessageSizeBytes)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		bytesMessage, err := s.conn.Read(messageBuff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.connCloseCh <- ErrConnectionClosed
				s.errCh <- err
				return
			}
			s.errCh <- errors.Join(err, ErrFailedToWaitMessage)
			continue
		}

		log.Printf("Received message from server, %d bytes", bytesMessage)

		exact := make([]byte, bytesMessage)
		copy(exact, messageBuff[:bytesMessage])

		s.messagesCh <- exact
	}
}

func (s *ServerSDK) CloseConnection() error {
	return s.conn.Close()
}

func (s *ServerSDK) WaitForClose() error {
	return <-s.connCloseCh
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

func (s *ServerSDK) PopMessage() (*protocol.RawMessage, error) {
	select {
	case message := <-s.messagesCh:
		return protocol.ParseRawMessage(message)

	case err := <-s.errCh:
		return nil, errors.Join(err, ErrFailedToWaitMessage)
	}
}
