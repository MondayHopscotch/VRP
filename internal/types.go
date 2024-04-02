package internal

import "math"

type Load struct {
	Number  int
	Pickup  Point
	Dropoff Point
}

type Point struct {
	X float64
	Y float64
}

func (p Point) DistanceTo(other Point) float64 {
	return math.Hypot(other.X-p.X, other.Y-p.Y)
}
