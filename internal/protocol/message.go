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
