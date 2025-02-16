package protocol

import (
	"encoding/binary"
	"errors"
	"log"
)

var (
	ErrFailedToEncodeMessage = errors.New("failed to encode message")
	ErrMessageTooShort       = errors.New("message is too short")
)

type RawMessage struct {
	Flags  byte
	Opcode uint32
	Data   []byte
}

type ParsedMessage[T any] struct {
	Flags  byte
	Opcode uint32
	Data   T
}

type MessageDecoder interface {
	Decode(data []byte) (any, error)
}

type MessageEncoder interface {
	Encode() ([]byte, error)
}

func (m *RawMessage) IsSuccess() bool {
	f := MessageFlags(m.Flags)
	return !f.HasFlag(MSG_FAIL_FLAG)
}

func (m *RawMessage) IsFailure() bool {
	f := MessageFlags(m.Flags)
	return f.HasFlag(MSG_FAIL_FLAG)
}

func BuildRawMessage(success bool, opcode uint32, payload MessageEncoder) ([]byte, error) {
	messageBuff := make([]byte, 5)

	flags := EmptyMessageFlags()
	if !success {
		flags.SetFlag(MSG_FAIL_FLAG)
	}
	messageBuff[0] = byte(flags)

	binary.BigEndian.PutUint32(messageBuff[1:5], opcode)

	if payload != nil {
		buff, err := payload.Encode()
		if err != nil {
			return nil, errors.Join(err, ErrFailedToEncodeMessage)
		}

		messageBuff = append(messageBuff, buff...)
	}
	log.Printf("Built message. [SIZE: %d bytes]", len(messageBuff))

	return messageBuff, nil
}

func ParseRawMessage(rawMessage []byte) (*RawMessage, error) {
	if len(rawMessage) < 5 {
		return nil, ErrMessageTooShort
	}

	flags := rawMessage[0]
	opcode := binary.BigEndian.Uint32(rawMessage[1:5])
	return &RawMessage{
		Flags:  flags,
		Opcode: opcode,
		Data:   rawMessage[5:],
	}, nil
}
