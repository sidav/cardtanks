package calc

/// This struct has two vectors: coordinates and direction.
type DirectedCoords struct {
	Coords IntVector2d
	Dir    IntVector2d
}

func NewDirectedCoords(x, y, dirX, dirY int) DirectedCoords {
	return DirectedCoords{
		Coords: NewIntVector2d(x, y),
		Dir: NewIntVector2d(dirX, dirY),
	}
}

func (dc *DirectedCoords) PointsToCoords() IntVector2d {
	return NewIntVector2d(dc.Coords.X + dc.Dir.X, dc.Coords.Y + dc.Dir.Y)
}

