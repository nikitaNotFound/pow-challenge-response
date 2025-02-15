package protocol

type Message struct {
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

func (m *Message) IsSuccess() bool {
	f := MessageFlags(m.Flags)
	return !f.HasFlag(MSG_FAIL_FLAG)
}

func (m *Message) IsFailure() bool {
	f := MessageFlags(m.Flags)
	return f.HasFlag(MSG_FAIL_FLAG)
}
