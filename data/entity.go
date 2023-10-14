package data

import "time"

type Entity struct {
	ID        int64     `bun:",pk,autoincrement" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
