package server_node

type ServerConfig struct {
	Address                 string
	MaxMessageSizeBytes     int
	ChallengeDifficulty     uint64
	MaxConnectionsPerClient int
}

func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		Address:                 "127.0.0.1:12345",
		MaxMessageSizeBytes:     1024,
		ChallengeDifficulty:     3,
		MaxConnectionsPerClient: 2,
	}
}
