package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
	"wordofwisdom/internal/pow"
	"wordofwisdom/internal/quotes"
)

const (
	PORT = "127.0.0.1:12345"
)

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func main() {
	time.Sleep(time.Second)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Client disconnected: %v", err)
			return
		}

		if msg.Type != "quote_request" {
			log.Printf("Invalid message type: %s", msg.Type)
			continue
		}

		challenge := pow.GenerateChallenge()
		if err := encoder.Encode(Message{Type: "challenge", Payload: challenge}); err != nil {
			log.Printf("Failed to send challenge: %v", err)
			return
		}

		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Failed to receive solution: %v", err)
			return
		}

		solution, ok := msg.Payload.(float64)
		if !ok {
			log.Printf("Invalid solution format from %s", conn.RemoteAddr())
			if err := encoder.Encode(Message{Type: "error", Payload: "Invalid solution format"}); err != nil {
				log.Printf("Failed to send error: %v", err)
			}
			continue
		}

		if !challenge.Verify(int(solution)) {
			log.Printf("Invalid solution from %s", conn.RemoteAddr())
			if err := encoder.Encode(Message{Type: "error", Payload: "Invalid solution"}); err != nil {
				log.Printf("Failed to send error: %v", err)
			}
			continue
		}

		quote := quotes.GetRandomQuote()
		if err := encoder.Encode(Message{Type: "quote", Payload: quote}); err != nil {
			log.Printf("Failed to send quote: %v", err)
			return
		}
	}
}
