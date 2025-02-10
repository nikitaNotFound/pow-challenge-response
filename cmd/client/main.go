package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"wordofwisdom/internal/pow"
)

const (
	SERVER_ADDR = "127.0.0.1:12345"
)

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func main() {
	conn, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	fmt.Println("Enter 'quote' to get a new quote or 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			return
		}

		if input == "quote" {
			// Request quote and solve challenge
			if err := getQuote(encoder, decoder); err != nil {
				log.Printf("Error getting quote: %v", err)
				continue
			}
		} else {
			fmt.Println("Invalid command. Enter 'quote' to get a new quote or 'exit' to quit")
		}
	}
}

func getQuote(encoder *json.Encoder, decoder *json.Decoder) error {
	if err := encoder.Encode(Message{Type: "quote_request"}); err != nil {
		return fmt.Errorf("failed to send quote request: %v", err)
	}

	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		return fmt.Errorf("failed to receive challenge: %v", err)
	}

	if msg.Type != "challenge" {
		return fmt.Errorf("expected challenge, got %s", msg.Type)
	}

	challengeData, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal challenge data: %v", err)
	}

	var challenge pow.Challenge
	if err := json.Unmarshal(challengeData, &challenge); err != nil {
		return fmt.Errorf("failed to unmarshal challenge: %v", err)
	}

	fmt.Println("Solving proof of work challenge...")
	started := time.Now()
	solution, err := challenge.Solve()
	if err != nil {
		return fmt.Errorf("failed to solve challenge: %v", err)
	}
	elapsed := time.Since(started)

	if err := encoder.Encode(Message{Type: "solution", Payload: solution}); err != nil {
		return fmt.Errorf("failed to send solution: %v", err)
	}

	if err := decoder.Decode(&msg); err != nil {
		return fmt.Errorf("failed to receive response: %v", err)
	}

	switch msg.Type {
	case "quote":
		fmt.Printf("\nChanllenge solution: %v\n", solution)
		fmt.Printf("\nTime taken to solve challenge: %v\n", elapsed)
		fmt.Printf("\nQuote: %v\n", msg.Payload)
	case "error":
		return fmt.Errorf("server error: %v", msg.Payload)

	default:
		return fmt.Errorf("unexpected response type: %s", msg.Type)
	}

	return nil
}
