package server_node

import (
	"context"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"wordofwisdom/pkg/protocol"
	"wordofwisdom/pkg/worker_pool"
)

type ServerHandler func(ctx *ServerContext) error

type TcpServer struct {
	maxMessageSizeBytes     int
	maxConnectionsPerClient int
	clientTimeout           time.Duration
	address                 string
	ctx                     context.Context

	handlers map[uint32]ServerHandler

	connections      map[string]int
	connectionsMutex sync.Mutex
	workerPool       *worker_pool.WorkerPool
}

func NewTcpServer(
	ctx context.Context,
	cfg *ServerConfig,
) *TcpServer {
	return &TcpServer{
		maxMessageSizeBytes:     cfg.MaxMessageSizeBytes,
		maxConnectionsPerClient: cfg.MaxConnectionsPerClient,
		clientTimeout:           time.Duration(cfg.ClientTimeoutMilliseconds) * time.Millisecond,
		address:                 cfg.Address,
		ctx:                     ctx,
		handlers:                make(map[uint32]ServerHandler),
		connections:             make(map[string]int),
		connectionsMutex:        sync.Mutex{},
		workerPool:              worker_pool.NewWorkerPool(cfg.WorkersAmount, ctx),
	}
}

func (s *TcpServer) RegisterHandler(
	opcode uint32,
	handler ServerHandler,
) {
	s.handlers[opcode] = handler
}

func (s *TcpServer) Run() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", s.address)

	s.workerPool.Start()

	for {
		select {
		case <-s.ctx.Done():
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
	if clientConnections >= s.maxConnectionsPerClient {
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

	serverCtx := NewServerContext(s.ctx, conn, s.maxMessageSizeBytes, s.clientTimeout)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		msg, err := serverCtx.WaitMessage()
		if err != nil {
			if errors.Is(err, ErrConnectionClosed) {
				log.Printf("Client %s disconnected", clientAddress)
				return
			}
			if errors.Is(err, ErrClientTimeout) {
				log.Printf("Client %s timed out", clientAddress)
				return
			}
			log.Printf("Failed to wait for message: %v", err)
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
