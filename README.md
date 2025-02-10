# Proof of Work Challenge-Response Server

A TCP server implementation protected from DDoS attacks using a Proof of Work (PoW) challenge-response protocol. After successfully solving the PoW challenge, the server sends a random quote from a collection.

## How It Works

1. **Client-Server Interaction**:
   - Client connects and requests a quote
   - Server sends a PoW challenge
   - Client solves the challenge (finds a nonce that produces a hash with required leading zeros)
   - Server verifies the solution and sends a quote if valid

2. **Proof of Work Implementation**:
   - Uses SHA-256 hashing algorithm
   - Challenge includes random data and timestamp
   - Solution requires finding a nonce that produces a hash with N leading zeros
   - Current difficulty is set to 6 leading zeros (can be adjusted in `internal/pow/pow.go`)
   - Higher difficulty values require more computational work:
     ```go
     // internal/pow/pow.go
     const (
         Difficulty = 6  // Increase this value for harder challenges
     )
     ```

3. **DDoS Protection**:
   - Each quote request requires significant computational work
   - Timestamp in challenge prevents replay attacks
   - CPU-intensive but memory-efficient
   - Easy to verify but time-consuming to solve

## Running the Project

### Prerequisites
- Go 1.23 or higher
- Make

### Windows
- Download and install [Go](https://go.dev/dl/)
- Download and install [Make](https://www.gnu.org/software/make/)
- Use `make build` to build the project
- Use `make run-all` to run the server and client

### Linux
- Install Go 1.23 or higher
- Use `make build` to build the project
- Use `make run-all` to run the server and client
