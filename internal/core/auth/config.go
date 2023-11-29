package auth

import "time"

// Config holds auth related configuration
type Config struct {
	TokenExp time.Duration // Token expiration time
	Secret   string        // Signing token secret
	Scopes   []string      // List of scopes required to access endpoint (default: none required)
}

func NewConfig(opts ...Option) *Config {
	c := new(Config)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(*Config)

func TokenExp(exp time.Duration) Option {
	return func(c *Config) {
		c.TokenExp = exp
	}
}

func Secret(secret string) Option {
	return func(c *Config) {
		c.Secret = secret
	}
}

func Scopes(scopes []string) Option {
	return func(c *Config) {
		c.Scopes = scopes
	}
}
