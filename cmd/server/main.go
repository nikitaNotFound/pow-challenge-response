package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
	"wordofwisdom/internal/pow"
)

const (
	PORT = "127.0.0.1:12345"
)

func main() {
	time.Sleep(time.Second)

	
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
