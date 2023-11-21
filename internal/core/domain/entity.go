package domain

import "time"

// Entity represents base for every persistent entity like User.
type Entity struct {
	ID        uint64    `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

// Page is generic struct that represents response made by page request.
type Page[T any] struct {
	TotalPages    int
	TotalElements int
	Elements      []T
}
