package main

import (
	"fmt"
	"github.com/fmiskovic/go-starter/internal/helper"
	"log/slog"
	"runtime"
	"strconv"
)

var defaultConfig ServerConfig

// ServerConfig holds server configuration.
type ServerConfig struct {
	ListenAddr   string
	DbConnString string
	MaxOpenConn  int
	MaxIdleConn  int
}

func init() {
	if err := helper.LoadEnvVars(); err != nil {
		slog.Warn("unable to locate .env file, default environment values will be used")
	}

	var (
		user   = helper.GetEnvOrDefault("DB_USER", "dbadmin")
		pass   = helper.GetEnvOrDefault("DB_PASSWORD", "dbadmin")
		host   = helper.GetEnvOrDefault("DB_HOST", "localhost:5432")
		dbName = helper.GetEnvOrDefault("DB_NAME", "go-db")

		conn = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbName)

		listenAddr = helper.GetEnvOrDefault("HTTP_LISTEN_ADDR", "localhost:8080")
	)

	numCpu := runtime.NumCPU() + 1
	// parsing DB_MAX_OPEN_CONN variable
	maxOpenConn, err := strconv.Atoi(helper.GetEnvOrDefault("DB_MAX_OPEN_CONN", strconv.Itoa(numCpu)))
	if err != nil {
		slog.Warn("error parsing DB_MAX_OPEN_CONN variable, using default", "error", err.Error())
		maxOpenConn = numCpu
	}
	// parsing DB_MAX_IDLE_CONN variable
	maxIdleConn, err := strconv.Atoi(helper.GetEnvOrDefault("DB_MAX_IDLE_CONN", strconv.Itoa(numCpu)))
	if err != nil {
		slog.Warn("error parsing DB_MAX_IDLE_CONN variable, using default", "error", err.Error())
		maxIdleConn = numCpu
	}

	defaultConfig = ServerConfig{
		ListenAddr:   listenAddr,
		DbConnString: conn,
		MaxOpenConn:  maxOpenConn,
		MaxIdleConn:  maxIdleConn,
	}
}
