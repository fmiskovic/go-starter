package data

import "time"

type Entity struct {
	ID        int64     `bun:",pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:"@created_at,notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:"@updated_at,notnull" json:"updatedAt"`
}
