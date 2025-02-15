package server

type ServerConfig struct {
	Address             string
	MaxMessageSizeBytes int
	ChallengeDifficulty int64
}

func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		Address:             "127.0.0.1:12345",
		MaxMessageSizeBytes: 1024,
		ChallengeDifficulty: 6,
	}
}
