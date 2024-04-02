package types

type Load struct {
	Number  int
	Pickup  Point
	Dropoff Point
}

// Cost returns the cost for the travel from Pickup to Dropoff for this load
func (l Load) Cost() float64 {
	return l.Pickup.DistanceTo(l.Dropoff)
}
