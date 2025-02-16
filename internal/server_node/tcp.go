package server_node

import (
	"context"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"wordofwisdom/pkg/protocol"
)

type ServerHandler func(ctx *ServerContext) error

type TcpServer struct {
	MaxMessageSizeBytes int
	Address             string
	Ctx                 context.Context

	handlers map[uint32]ServerHandler
}

func NewTcpServer(ctx context.Context, address string, maxMessageSizeBits int) *TcpServer {
	return &TcpServer{
		MaxMessageSizeBytes: maxMessageSizeBits,
		Address:             address,
		Ctx:                 ctx,
		handlers:            make(map[uint32]ServerHandler),
	}
}

func (s *TcpServer) RegisterHandler(
	opcode uint32,
	handler ServerHandler,
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

	clientIp := conn.RemoteAddr().String()
	log.Printf("New connection established with ip: %s", clientIp)

	serverCtx := NewServerContext(s.Ctx, conn, s.MaxMessageSizeBytes)

	for {
		select {
		case <-s.Ctx.Done():
			return
		default:
		}

		log.Printf("Waiting for message from client: %s", clientIp)

		msg, err := serverCtx.WaitMessage()
		if err != nil {
			log.Printf("Failed to wait for message: %v", err)
			if errors.Is(err, ErrConnectionClosed) {
				log.Printf("Client %s disconnected", clientIp)
				return
			}
			continue
		}

		handler, ok := s.handlers[msg.Opcode]
		if !ok {
			log.Printf("No handler found for opcode: %d", msg.Opcode)
			msgBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(msgBuf, protocol.ERR_CODE_INVALID_OPCODE)
			conn.Write(msgBuf)
			continue
		}
		handler(serverCtx)
	}
}
