package v1

import (
	"fmt"
	"time"
)

type RouterConfig struct {
	EnableRateLimit   bool
	RateLimitRequests int
	RateLimitWindow   time.Duration
	EnableMetrics     bool
	EnableCORS        bool
	AllowedOrigins    []string
	ServiceConfig     ServiceConfig
}

func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		EnableRateLimit:   true,
		RateLimitRequests: 100,
		RateLimitWindow:   time.Minute,
		EnableMetrics:     true,
		EnableCORS:        false,
		AllowedOrigins:    []string{},
		ServiceConfig: ServiceConfig{
			PublicKeyLength: PublicKeyBytes,
			SecretKeyLength: SecretKeyBytes,
			Timeout:         DefaultTimeout,
		},
	}
}

func ProductionRouterConfig() RouterConfig {
	return RouterConfig{
		EnableRateLimit:   true,
		RateLimitRequests: 60,
		RateLimitWindow:   time.Minute,
		EnableMetrics:     true,
		EnableCORS:        true,
		AllowedOrigins:    []string{},
		ServiceConfig: ServiceConfig{
			PublicKeyLength: PublicKeyBytes,
			SecretKeyLength: SecretKeyBytes,
			Timeout:         10 * time.Second,
		},
	}
}

func DevelopmentRouterConfig() RouterConfig {
	return RouterConfig{
		EnableRateLimit:   false,
		RateLimitRequests: 1000,
		RateLimitWindow:   time.Minute,
		EnableMetrics:     true,
		EnableCORS:        true,
		AllowedOrigins:    []string{"*"},
		ServiceConfig: ServiceConfig{
			PublicKeyLength: PublicKeyBytes,
			SecretKeyLength: SecretKeyBytes,
			Timeout:         30 * time.Second,
		},
	}
}
func (c RouterConfig) Validate() error {
	if c.EnableRateLimit {
		if c.RateLimitRequests <= 0 {
			return fmt.Errorf("rate limit requests must be positive")
		}
		if c.RateLimitWindow <= 0 {
			return fmt.Errorf("rate limit window must be positive")
		}
	}

	if err := c.ServiceConfig.Validate(); err != nil {
		return fmt.Errorf("invalid service config: %w", err)
	}

	return nil
}
