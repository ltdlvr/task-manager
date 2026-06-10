package model

import "time"

type Task struct {
	ID          uint64
	ColumnID    uint64
	Title       string
	Description string
	Position    int
	CreatedAt   time.Time
}

func (t *Task) GetPosition() int {
	return t.Position
}
