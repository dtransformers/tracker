package constants

type AppDefaultConfig struct {
	Port int
}

var DefaultConfig = AppDefaultConfig{
	Port: 3000,
}

func GetDefaultConfig() AppDefaultConfig {
	return DefaultConfig
}
