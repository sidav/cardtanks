package calc

type IntVector2d struct {
	X, Y int
}

func NewIntVector2d(x, y int) IntVector2d {
	return IntVector2d{X: x, Y: y}
}

func (v *IntVector2d) Is(x, y int) bool {
	return v.X == x && v.Y == y
}

func (v *IntVector2d) SetEqTo(x, y int) {
	v.X = x
	v.Y = y
}

// func (v *IntVector2d) AdjacentTo(v2 IntVector2d) bool {
// 	return v.X == v2.X && (v.Y == v2.Y-1 || v.Y == v2.Y+1) ||
// 		v.Y == v2.Y && (v.X == v2.X-1 || v.X == v2.X+1)
// }

func (v *IntVector2d) Unwrap() (int, int) {
	return v.X, v.Y
}

func (v *IntVector2d) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v *IntVector2d) Normalize() {
	if v.X != 0 {
		v.X = v.X / IntAbs(v.X)
	}
	if v.Y != 0 {
		v.Y = v.Y / IntAbs(v.Y)
	}
}

func (v *IntVector2d) RotateCCW() {
	v.X, v.Y = v.Y, -v.X
}

func (v *IntVector2d) RotateCW() {
	v.X, v.Y = -v.Y, v.X
}
