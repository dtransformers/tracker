package constants

import "time"

type AppDefaultConfig struct {
	Port int
}

var DefaultConfig = AppDefaultConfig{
	Port: 3000,
}

func GetDefaultConfig() AppDefaultConfig {
	return DefaultConfig
}

// DefaultServiceTimeout is the standard context timeout applied in service methods.
const DefaultServiceTimeout = 5 * time.Second
