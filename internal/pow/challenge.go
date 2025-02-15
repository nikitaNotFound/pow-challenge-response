package pow

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Challenge struct {
	Data           [16]byte
	Timestamp      int64
	Difficulty     int64
	ExpectedPrefix []byte
}

func GenerateChallenge(difficulty int64) *Challenge {
	return &Challenge{
		Data:           generateRandom16xHash(),
		Timestamp:      time.Now().Unix(),
		Difficulty:     difficulty,
		ExpectedPrefix: generateExpectedPrefix(difficulty),
	}
}

func NewChallenge(data [16]byte, timestamp int64, difficulty int64) *Challenge {
	return &Challenge{
		Data:           data,
		Timestamp:      timestamp,
		Difficulty:     difficulty,
		ExpectedPrefix: generateExpectedPrefix(difficulty),
	}
}

func generateExpectedPrefix(difficulty int64) []byte {
	return []byte(strings.Repeat("0", int(difficulty)))
}

func (c *Challenge) Solve() (int64, error) {
	target := c.ExpectedPrefix
	nonce := int64(0)

	for {
		hash := calculateHash(c.Data, c.Timestamp, nonce)
		if bytes.Equal(hash[:len(target)], target) {
			return nonce, nil
		}
		nonce++
	}
}

func (c *Challenge) Verify(nonce int64) bool {
	hash := calculateHash(c.Data, c.Timestamp, nonce)
	return bytes.Equal(hash[:len(c.ExpectedPrefix)], c.ExpectedPrefix)
}

func calculateHash(data [16]byte, timestamp int64, nonce int64) [32]byte {
	input := fmt.Sprintf("%s%d%d", data, timestamp, nonce)
	return sha256.Sum256([]byte(input))
}

func generateRandom16xHash() [16]byte {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b [16]byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return b
}
