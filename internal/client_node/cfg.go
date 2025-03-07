package client_node

type ClientConfig struct {
	ServerAddress       string
	MaxMessageSizeBytes int
	PopMessageTimeoutMs int
}

func GetClientConfig() *ClientConfig {
	return &ClientConfig{
		ServerAddress:       "127.0.0.1:12345",
		MaxMessageSizeBytes: 1024,
		PopMessageTimeoutMs: 15000,
	}
}
