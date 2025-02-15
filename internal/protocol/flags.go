package protocol

type MessageFlags byte

func EmptyMessageFlags() MessageFlags {
	return 0x0
}

// Message flags
// Only first flag is used for now to identify success/failure of message.
// Other flags are reserved for future use.
// All operations for operating flags are implemented using bitwise operations.
const (
	MSG_FAIL_FLAG MessageFlags = 1 << iota // 00000001
	FLAG_2                                 // 00000010
	FLAG_3                                 // 00000100
	FLAG_4                                 // 00001000
	FLAG_5                                 // 00010000
	FLAG_6                                 // 00100000
	FLAG_7                                 // 01000000
	FLAG_8                                 // 10000000
)

func (f *MessageFlags) SetFlag(flag MessageFlags) {
	*f = *f | flag
}

func (f *MessageFlags) ClearFlag(flag MessageFlags) {
	*f = *f &^ flag
}

func (f *MessageFlags) HasFlag(flag MessageFlags) bool {
	return *f&flag != 0
}
