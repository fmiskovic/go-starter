package database

import (
	"log/slog"
	"time"

	"github.com/fmiskovic/go-starter/util"
	_ "github.com/lib/pq"
)

const ConnString = "postgresql://%s:%s@%s/%s?sslmode=disable"

var DefaultConfig Config

type Config struct {
	User     string
	Password string
	Host     string
	Name     string
	Timeout  time.Duration
}

func init() {
	if err := util.LoadEnvVars(); err != nil {
		slog.Warn("unable to locate .env file, default environment values will be used")
	}

	DefaultConfig = Config{
		User:     util.GetEnvOrDefault("DB_USER", "dbadmin"),
		Password: util.GetEnvOrDefault("DB_PASSWORD", "dbadmin"),
		Host:     util.GetEnvOrDefault("DB_HOST", "localhost:5432"),
		Name:     util.GetEnvOrDefault("DB_NAME", "go-database"),
		Timeout:  time.Second * 10,
	}
}

//func init() {
//	if err := util.LoadEnvVars(); err != nil {
//		slog.Warn("unable to locate .env file, default environment values will be used")
//	}
//
//	var (
//		user     = util.GetEnvOrDefault("DB_USER", "dbadmin")
//		password = util.GetEnvOrDefault("DB_PASSWORD", "dbadmin")
//		host     = util.GetEnvOrDefault("DB_HOST", "localhost:5432")
//		name     = util.GetEnvOrDefault("DB_NAME", "go-database")
//		uri      = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)
//		sqldb    = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
//		once     = sync.Once{}
//	)
//	fmt.Println(uri)
//	once.Do(func() {
//		Bun = bun.NewDB(sqldb, pgdialect.New())
//		if util.IsDev() {
//			Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
//		}
//	})
//}
