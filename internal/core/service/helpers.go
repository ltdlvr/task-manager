package service

import (
	"fmt"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
)

type positioned interface {
	GetPosition() int
}

func getNeighbours[T positioned](items []T, targetPos int) (prev, next *int, err error) {
	if targetPos < 0 || targetPos > len(items) {
		return nil, nil, fmt.Errorf("targetPos %d out of range [0, %d]: %w", targetPos, len(items), db.ErrOutOfRange)
	}
	if targetPos > 0 {
		pos := items[targetPos-1].GetPosition()
		prev = &pos
	}
	if targetPos < len(items) {
		pos := items[targetPos].GetPosition()
		next = &pos
	}
	return
}
