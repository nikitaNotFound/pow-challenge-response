package server_node

import (
	"context"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"wordofwisdom/pkg/protocol"
	"wordofwisdom/pkg/worker_pool"
)

type ServerHandler func(ctx *ServerContext) error

type TcpServer struct {
	MaxMessageSizeBytes     int
	MaxConnectionsPerClient int
	Address                 string
	Ctx                     context.Context

	handlers map[uint32]ServerHandler

	connections      map[string]int
	connectionsMutex sync.Mutex
	workerPool       *worker_pool.WorkerPool
}

func NewTcpServer(
	ctx context.Context,
	address string,
	maxMessageSizeBits int,
	maxConnectionsPerClient int,
	workersAmount int,
) *TcpServer {
	return &TcpServer{
		MaxMessageSizeBytes:     maxMessageSizeBits,
		MaxConnectionsPerClient: maxConnectionsPerClient,
		Address:                 address,
		Ctx:                     ctx,
		handlers:                make(map[uint32]ServerHandler),
		connections:             make(map[string]int),
		connectionsMutex:        sync.Mutex{},
		workerPool:              worker_pool.NewWorkerPool(workersAmount, ctx),
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

	s.workerPool.Start()

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

		s.workerPool.RunWork(func() {
			s.handleNewConnection(conn)
		})
	}
}

func (s *TcpServer) reserveClientConnection(clientIp string) error {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()
	clientConnections := s.connections[clientIp]
	if clientConnections >= s.MaxConnectionsPerClient {
		log.Printf("Max connections per client reached for ip: %s, current connections: %d", clientIp, clientConnections)
		return errors.New("max connections per client reached")
	}
	s.connections[clientIp]++
	return nil
}

func (s *TcpServer) releaseClientConnection(clientIp string) {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()
	s.connections[clientIp]--
}

func (s *TcpServer) handleNewConnection(conn net.Conn) {
	clientIp := strings.Split(conn.RemoteAddr().String(), ":")[0]

	if err := s.reserveClientConnection(clientIp); err != nil {
		conn.Close()
		return
	}

	clientAddress := conn.RemoteAddr().String()
	log.Printf("New connection established with ip: %s", clientAddress)

	defer func() {
		conn.Close()
		s.releaseClientConnection(clientIp)
	}()

	serverCtx := NewServerContext(s.Ctx, conn, s.MaxMessageSizeBytes)

	for {
		select {
		case <-s.Ctx.Done():
			return
		default:
		}

		log.Printf("Waiting for message from client: %s", clientAddress)

		msg, err := serverCtx.WaitMessage()
		if err != nil {
			log.Printf("Failed to wait for message: %v", err)
			if errors.Is(err, ErrConnectionClosed) {
				log.Printf("Client %s disconnected", clientAddress)
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
