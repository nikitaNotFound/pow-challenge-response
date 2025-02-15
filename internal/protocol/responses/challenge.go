package responses

import (
	"encoding/binary"
	"errors"
)

type ChallengeResponse struct {
	Data           [16]byte
	Timestamp      uint64
	Difficulty     uint64
	ExpectedPrefix []byte
}

func (cr *ChallengeResponse) Encode() ([]byte, error) {
	buff := make([]byte, 16+8+8+len(cr.ExpectedPrefix))
	copy(buff[:16], cr.Data[:])
	binary.BigEndian.PutUint64(buff[16:24], uint64(cr.Timestamp))
	binary.BigEndian.PutUint64(buff[24:32], uint64(cr.Difficulty))
	copy(buff[32:], cr.ExpectedPrefix)
	return buff, nil
}

func (cr *ChallengeResponse) Decode(buff []byte) error {
	if len(buff) <= 16+8+8 {
		return errors.New("invalid challenge response")
	}

	var dataBuff [16]byte
	copy(dataBuff[:], buff[:16])

	timestamp := binary.BigEndian.Uint64(buff[16:24])
	difficulty := binary.BigEndian.Uint64(buff[24:32])

	if len(buff)-32 != int(difficulty) {
		return errors.New("expected prefix length not the same as difficulty")
	}

	expectedPrefixBuff := make([]byte, difficulty)
	copy(expectedPrefixBuff, buff[32:])

	cr.Data = dataBuff
	cr.Timestamp = timestamp
	cr.Difficulty = difficulty
	cr.ExpectedPrefix = expectedPrefixBuff

	return nil
}
