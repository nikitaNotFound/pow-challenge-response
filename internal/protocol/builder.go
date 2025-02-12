package protocol

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}
