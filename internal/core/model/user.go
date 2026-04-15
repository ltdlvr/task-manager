package model

import "time"

type User struct {
	ID        uint64
	Name      string
	Password  string
	CreatedAt time.Time
}
