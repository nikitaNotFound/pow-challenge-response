package requests

import (
	"encoding/binary"
	"errors"
)

type ChallengeProofRequest struct {
	Nonce int64
}

func (cpr *ChallengeProofRequest) Encode() ([]byte, error) {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(cpr.Nonce))
	return buff, nil
}

func (cpr *ChallengeProofRequest) Decode(buff []byte) error {
	if len(buff) != 8 {
		return errors.New("invalid challenge proof request")
	}

	cpr.Nonce = int64(binary.BigEndian.Uint64(buff))
	return nil
}
