package ports

import (
	"github.com/uptrace/bun"
)

type Db interface {
	OpenDb() (*bun.DB, error)
}
