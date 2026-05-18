package model

import "time"

type Column struct {
	ID        uint64
	BoardID   uint64
	Name      string
	Position  int
	CreatedAt time.Time
}
