package server

import (
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"
	"wordofwisdom/internal/messages"
)

type TcpServer struct {
	MaxMessageSizeBytes int
	Address             string
	Ctx                 context.Context

	handlers map[uint32]func(context.Context, io.Writer, io.Reader, []byte) error
}

func NewTcpServer(ctx context.Context, address string, maxMessageSizeBits int) *TcpServer {
	return &TcpServer{
		MaxMessageSizeBytes: maxMessageSizeBits,
		Address:             address,
		Ctx:                 ctx,
	}
}

func (s *TcpServer) RegisterHandler(
	opcode uint32,
	handler func(context.Context, io.Writer, io.Reader, []byte) error,
) {
	s.handlers[opcode] = handler
}

func (s *TcpServer) Run() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", s.Address)

	for {
		select {
		case <-s.Ctx.Done():
			return nil
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleNewConnection(conn)
	}
}

func (s *TcpServer) handleNewConnection(conn net.Conn) {
	defer conn.Close()

	for {
		select {
		case <-s.Ctx.Done():
			return
		default:
		}

		log.Printf("Waiting for message...")

		messageBuff := make([]byte, s.MaxMessageSizeBytes)
		bytesMessage, err := conn.Read(messageBuff)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client")
				return
			}
			log.Printf("Failed to read message: %v", err)
			continue
		}

		opcode := binary.BigEndian.Uint32(messageBuff[:4])
		handler, ok := s.handlers[opcode]
		if !ok {
			log.Printf("No handler found for opcode: %d", opcode)
			msgBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(msgBuf, messages.ERR_CODE_INVALID_OPCODE)
			conn.Write(msgBuf)
			continue
		}
		handler(s.Ctx, conn, conn, messageBuff[4:bytesMessage])
	}
}
