package internal

import (
	"encoding/json"
	"fmt"
	"math"
)

type Route []Load

func (r Route) PrintLoadNumbers() {
	nums := make([]int, len(r))
	for i, l := range r {
		nums[i] = l.Number
	}
	out, _ := json.Marshal(nums)
	fmt.Println(string(out))
}

// DistanceTo returns the distance from the end of the route to the given point
func (r Route) DistanceTo(p Point) float64 {
	if len(r) == 0 {
		return Point{X: 0, Y: 0}.DistanceTo(p)
	}

	return r[len(r)-1].Dropoff.DistanceTo(p)
}

// CurrentCost returns the cost of the route up to the last destination currently in the list
func (r Route) CurrentCost() float64 {
	// TODO: need to calculate total cost somewhere else. For the sake of searching, we care about time, not cost
	cost := 0.0
	p := Point{X: 0, Y: 0}
	for _, l := range r {
		cost += p.DistanceTo(l.Pickup)
		cost += l.Cost()
		p = l.Dropoff
	}

	return cost
}

// CurrentCost returns the current cost of the route and driver, including returning to the hub
func (r Route) CurrentCompletionCost() float64 {
	cost := r.CurrentCost()
	return cost + r.DistanceTo(Point{X: 0, Y: 0})
}

func (r Route) CostWithLoad(load Load) float64 {
	return r.CurrentCost() + r.DistanceTo(load.Pickup) + load.Cost()
}

// CurrentCost returns the current cost of the route and driver, including returning to the hub
func (r Route) CompletionCostWithLoad(load Load) float64 {
	cost := r.CostWithLoad(load)
	return cost + load.Dropoff.DistanceTo(Point{X: 0, Y: 0})
}

// CostIncreaseWithLoad returns the increase of cost, assuming the route will return to hub after the provided load is complete
func (r Route) CostIncreaseWithLoad(load Load) float64 {
	return r.DistanceTo(load.Pickup) + load.Cost() + load.Dropoff.DistanceTo(Point{})
}

type Load struct {
	Number  int
	Pickup  Point
	Dropoff Point
}

// Cost returns the cost for the travel from Pickup to Dropoff for this load
func (l Load) Cost() float64 {
	return l.Pickup.DistanceTo(l.Dropoff)
}

type Point struct {
	X float64
	Y float64
}

func (p Point) DistanceTo(other Point) float64 {
	return math.Hypot(other.X-p.X, other.Y-p.Y)
}
