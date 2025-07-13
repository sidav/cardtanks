package calc

import "strconv"

func IntAbs(x int) int {
	y := x >> (strconv.IntSize - 1)
	return (x ^ y) - y
}

func ApproxDistanceInt(x1, y1, x2, y2 int) int {
	diffX := IntAbs(x1 - x2)
	diffY := IntAbs(y1 - y2)
	if diffX > diffY {
		return diffX + (diffY >> 1)
	}
	return diffY + (diffX >> 1)
}
