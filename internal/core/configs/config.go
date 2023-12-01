package configs

import (
	"time"
)

// Config holds auth related configuration
type AuthConfig struct {
	TokenExp time.Duration // Token expiration time
	Secret   string        // Signing token secret
	Scopes   []string      // List of scopes required to access endpoint (default: none required)
}

func NewAuthConfig(opts ...AuthConfigOptions) AuthConfig {
	cfg := &AuthConfig{Secret: "secret"}
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

type AuthConfigOptions func(*AuthConfig)

func TokenExp(exp time.Duration) AuthConfigOptions {
	return func(ac *AuthConfig) {
		ac.TokenExp = exp
	}
}

func Secret(s string) AuthConfigOptions {
	return func(ac *AuthConfig) {
		ac.Secret = s
	}
}

func Scopes(s []string) AuthConfigOptions {
	return func(ac *AuthConfig) {
		ac.Scopes = s
	}
}
