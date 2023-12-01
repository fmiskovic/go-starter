package main

import (
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"time"

	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/utils"
)

var defaultConfig ServerConfig

// ServerConfig holds server configuration.
type ServerConfig struct {
	ListenAddr   string
	DbConnString string
	MaxOpenConn  int
	MaxIdleConn  int
	AuthConfig   configs.AuthConfig
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

func initDefaultAuthConfig() configs.AuthConfig {
	// parsing DB_MAX_IDLE_CONN variable
	expTime, err := strconv.Atoi(utils.GetEnvOrDefault("AUTH_JWT_EXP_TIME", "24"))
	if err != nil {
		slog.Warn("error parsing AUTH_JWT_EXP_TIME variable, using default", "error", err.Error())
		expTime = 24
	}

	secret := utils.GetEnvOrDefault("AUTH_JWT_SECRET", "secret")

	slog.Info("default auth config is initialized")
	return *&configs.AuthConfig{
		TokenExp: time.Duration(expTime),
		Secret:   secret,
	}
}
