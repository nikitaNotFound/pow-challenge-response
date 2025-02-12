package server

type ServerConfig struct {
	Address             string
	MaxMessageSizeBytes int
	ChallengeDifficulty int
}

func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		Address:             "localhost:8080",
		MaxMessageSizeBytes: 1024,
		ChallengeDifficulty: 6,
	}
}
