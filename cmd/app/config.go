package main

import (
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"time"

	"github.com/fmiskovic/go-starter/internal/utils"
)

var defaultConfig ServerConfig

// ServerConfig holds server configuration.
type ServerConfig struct {
	ListenAddr   string
	DbConnString string
	MaxOpenConn  int
	MaxIdleConn  int
	AuthConfig   AuthConfig
}

// AuthConfig holds auth related configuration
type AuthConfig struct {
	TokenExp time.Time // Token expiration time in hours
	Secret   string    // Signing token secret
	Scopes   []string  // List of scopes required to access endpoint (default: none required)
	Enabled  bool      // Auth enable/disable flag - default disabled
}

func init() {
	if err := utils.LoadEnvVars(); err != nil {
		slog.Warn("unable to locate .env file, default environment values will be used")
	}

	defaultConfig = initDefaultServerConfig()
}

func initDefaultServerConfig() ServerConfig {
	var (
		user   = utils.GetEnvOrDefault("DB_USER", "dbadmin")
		pass   = utils.GetEnvOrDefault("DB_PASSWORD", "dbadmin")
		host   = utils.GetEnvOrDefault("DB_HOST", "localhost:5432")
		dbName = utils.GetEnvOrDefault("DB_NAME", "go-db")

		conn = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbName)

		listenAddr = utils.GetEnvOrDefault("HTTP_LISTEN_ADDR", ":8080")
	)

	numCpu := runtime.NumCPU() + 1
	// parsing DB_MAX_OPEN_CONN variable
	maxOpenConn, err := strconv.Atoi(utils.GetEnvOrDefault("DB_MAX_OPEN_CONN", strconv.Itoa(numCpu)))
	if err != nil {
		slog.Warn("error parsing DB_MAX_OPEN_CONN variable, using default", "error", err.Error())
		maxOpenConn = numCpu
	}
	// parsing DB_MAX_IDLE_CONN variable
	maxIdleConn, err := strconv.Atoi(utils.GetEnvOrDefault("DB_MAX_IDLE_CONN", strconv.Itoa(numCpu)))
	if err != nil {
		slog.Warn("error parsing DB_MAX_IDLE_CONN variable, using default", "error", err.Error())
		maxIdleConn = numCpu
	}

	slog.Info("default server config is initialized")

	return ServerConfig{
		ListenAddr:   listenAddr,
		DbConnString: conn,
		MaxOpenConn:  maxOpenConn,
		MaxIdleConn:  maxIdleConn,
		AuthConfig:   initDefaultAuthConfig(),
	}
}

func initDefaultAuthConfig() AuthConfig {
	enabled, err := strconv.ParseBool(utils.GetEnvOrDefault("AUTH_ENABLED", "false"))
	if err != nil {
		slog.Warn("error parsing AUTH_ENABLED variable, auth will not be initialized", "error", err.Error())
	}

	if !enabled {
		slog.Info("auth is disabled")
		return AuthConfig{Enabled: false}
	}

	slog.Info("default auth config is initialized")
	return AuthConfig{}
}
