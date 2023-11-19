package domain

import "time"

type Entity struct {
	ID        uint64    `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
