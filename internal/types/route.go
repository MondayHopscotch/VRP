package types

import (
	"encoding/json"
	"fmt"
)

// DriverCost is the base cost of each driver to drive a Route
const DriverCost = 500.0

// DriverMaxTime is the maximum shift length each individual driver is permitted
const DriverMaxTime = 12.0 * 60

// Route is an alias to allow easier calculations to be done on a slice of Loads
type Route []Load

// GetTotalCostOfRoutes provides a shorthand for totalling the cost of a slice of routes
func GetTotalCostOfRoutes(routes []Route) float64 {
	totalCost := 0.0

	for _, r := range routes {
		totalCost += r.TotalCostWithDriver()
	}

	return totalCost
}

// PrintLoadNumbers prints routing arrays to console. Useful for testing output by external utilities
func (r Route) PrintLoadNumbers() {
	nums := make([]int, len(r))
	for i, l := range r {
		nums[i] = l.Number
	}
	out, _ := json.Marshal(nums)
	fmt.Println(string(out))
}

// DistanceTo returns the distance from the last dropoff of the route to the provided point
func (r Route) DistanceTo(p Point) float64 {
	if len(r) == 0 {
		return Point{X: 0, Y: 0}.DistanceTo(p)
	}

	return r[len(r)-1].Dropoff.DistanceTo(p)
}

// CurrentTime returns the time of the route up to the last load dropoff
func (r Route) CurrentTime() float64 {
	time := 0.0
	p := Point{X: 0, Y: 0}
	for _, l := range r {
		time += p.DistanceTo(l.Pickup)
		time += l.Cost()
		p = l.Dropoff
	}

	return time
}

// CurrentCompletionTime returns the current time of the route, including returning to the hub
func (r Route) CurrentCompletionTime() float64 {
	time := r.CurrentTime()
	return time + r.DistanceTo(Point{X: 0, Y: 0})
}

// TotalCostWithDriver returns the total cost of the route including the driver base cost
func (r Route) TotalCostWithDriver() float64 {
	return r.CurrentCompletionTime() + DriverCost
}

// TimeWithLoad returns the time of the route if the provided load were added to it
func (r Route) TimeWithLoad(load Load) float64 {
	return r.CurrentTime() + r.DistanceTo(load.Pickup) + load.Cost()
}

// CompletionTimeWithLoad returns the current time of the route, including returning to the hub
func (r Route) CompletionTimeWithLoad(load Load) float64 {
	cost := r.TimeWithLoad(load)
	return cost + load.Dropoff.DistanceTo(Point{X: 0, Y: 0})
}

// CompletionTimeIncreaseWithLoad returns the increase of time, assuming the route will return to hub after the provided load is complete
func (r Route) CompletionTimeIncreaseWithLoad(load Load) float64 {
	return r.DistanceTo(load.Pickup) + load.Cost() + load.Dropoff.DistanceTo(HubPoint()) - r.DistanceTo(HubPoint())
}
