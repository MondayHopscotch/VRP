package types

import (
	"math"
)

type Point struct {
	X float64
	Y float64
}

func HubPoint() Point {
	return Point{X: 0, Y: 0}
}

func (p Point) DistanceTo(other Point) float64 {
	return math.Hypot(other.X-p.X, other.Y-p.Y)
}
