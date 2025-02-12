package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Challenge struct {
	Data       string
	Timestamp  int64
	Difficulty int
}

func GenerateChallenge(difficulty int) *Challenge {
	return &Challenge{
		Data:       generateRandomString(16),
		Timestamp:  time.Now().Unix(),
		Difficulty: difficulty,
	}
}

func NewChallenge(data string, timestamp int64, difficulty int) *Challenge {
	return &Challenge{
		Data:       data,
		Timestamp:  timestamp,
		Difficulty: difficulty,
	}
}

func (c *Challenge) Solve() (int, error) {
	target := strings.Repeat("0", c.Difficulty)
	nonce := 0

	for {
		hash := calculateHash(c.Data, c.Timestamp, nonce)
		if strings.HasPrefix(hash, target) {
			return nonce, nil
		}
		nonce++
	}
}

func (c *Challenge) Verify(nonce int) bool {
	hash := calculateHash(c.Data, c.Timestamp, nonce)
	target := strings.Repeat("0", c.Difficulty)
	return strings.HasPrefix(hash, target)
}

func calculateHash(data string, timestamp int64, nonce int) string {
	input := fmt.Sprintf("%s%d%d", data, timestamp, nonce)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
