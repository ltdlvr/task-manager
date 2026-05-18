package tool //NOTE - Возможно, стоит разнести в разные пакеты тулзы

import (
	"fmt"
)

const (
	DefaultStep = 1024
	Startpos    = 1024
)

func CalculateNewPosition(prev, next *int) (int, error) {
	switch {
	case prev == nil && next == nil:
		return Startpos, nil
	case prev == nil:
		return *next - DefaultStep, nil
	case next == nil:
		return *prev + DefaultStep, nil
	default:
		mid := (*prev + *next) / 2
		if mid == *prev || mid == *next {
			return 0, fmt.Errorf("no space between %d and %d, rebalance needed", prev, *next)
		}
		return mid, nil
	}
}
