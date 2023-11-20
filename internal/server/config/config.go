package config

import (
	"fmt"
	"github.com/fmiskovic/go-starter/internal/util"
	"log/slog"
	"runtime"
	"strconv"
)

var DefaultConfig ServerConfig

// ServerConfig holds server configuration.
type ServerConfig struct {
	ListenAddr   string
	DbConnString string
	MaxOpenConn  int
	MaxIdleConn  int
}

func init() {
	if err := util.LoadEnvVars(); err != nil {
		slog.Warn("unable to locate .env file, default environment values will be used")
	}

	var (
		user   = util.GetEnvOrDefault("DB_USER", "dbadmin")
		pass   = util.GetEnvOrDefault("DB_PASSWORD", "dbadmin")
		host   = util.GetEnvOrDefault("DB_HOST", "localhost:5432")
		dbName = util.GetEnvOrDefault("DB_NAME", "go-db")

		conn = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbName)

		listenAddr = util.GetEnvOrDefault("HTTP_LISTEN_ADDR", "localhost:8080")
	)

	numCpu := runtime.NumCPU() + 1
	// parsing DB_MAX_OPEN_CONN variable
	maxOpenConn, err := strconv.Atoi(util.GetEnvOrDefault("DB_MAX_OPEN_CONN", strconv.Itoa(numCpu)))
	if err != nil {
		slog.Warn("error parsing DB_MAX_OPEN_CONN variable, using default", "error", err.Error())
		maxOpenConn = numCpu
	}
	// parsing DB_MAX_IDLE_CONN variable
	maxIdleConn, err := strconv.Atoi(util.GetEnvOrDefault("DB_MAX_IDLE_CONN", strconv.Itoa(numCpu)))
	if err != nil {
		slog.Warn("error parsing DB_MAX_IDLE_CONN variable, using default", "error", err.Error())
		maxIdleConn = numCpu
	}

	DefaultConfig = New(
		WithListenAddr(listenAddr),
		WithDbConnString(conn),
		WithMaxOpenConn(maxOpenConn),
		WithIdleOpenConn(maxIdleConn),
	)
}

// New instantiate new ServerConfig with the given options.
func New(opts ...ServerConfigOption) ServerConfig {
	conf := &ServerConfig{}
	for _, opt := range opts {
		opt(conf)
	}
	return *conf
}

type ServerConfigOption func(*ServerConfig)

func WithListenAddr(addr string) ServerConfigOption {
	return func(c *ServerConfig) {
		c.ListenAddr = addr
	}
}

func WithDbConnString(conn string) ServerConfigOption {
	return func(c *ServerConfig) {
		c.DbConnString = conn
	}
}

func WithMaxOpenConn(maxOpen int) ServerConfigOption {
	return func(c *ServerConfig) {
		c.MaxOpenConn = maxOpen
	}
}

func WithIdleOpenConn(maxIdle int) ServerConfigOption {
	return func(c *ServerConfig) {
		c.MaxIdleConn = maxIdle
	}
}
