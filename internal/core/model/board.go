package model

import "time"

type Board struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
}
