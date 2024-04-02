package routing

import (
	"fmt"
	"math"
	"vrp/internal"
)

type Solver interface {
	PlanRoutes(driverCount int) []internal.Route
}

func NewSolver(loads []internal.Load) Solver {
	return &BasicSolver{loads: loads}
}

// BasicSolver plans routes basedon a simple greedy algorithm for routing our drives
type BasicSolver struct {
	loads []internal.Load
}

func (b *BasicSolver) PlanRoutes(driverCount int) []internal.Route {
	// TODO: require at least one driver
	// TODO: Need to return error if we can't find a set of routes

	// Add all drivers to our pool
	// Loop:
	//  Find closest load to each driver that *doesn't put them into overtime with the trip to hub taken into account*
	//      Note that this may not be the closest load, as something further away, but also on the way to the hub may be a doable option
	//  Assign highest priority load to the appropriate driver
	//
	// We have to be careful not to strand drivers

	routes := make([]internal.Route, driverCount)

	remainingLoads := make(map[int]internal.Load)
	for _, l := range b.loads {
		remainingLoads[l.Number] = l
	}

	iter := 0

	driver := 0
	load := 0
	loadCost := math.MaxFloat64
	for len(remainingLoads) > 0 {
		iter++
		driver = 0
		load = 0
		loadCost = math.MaxFloat64

		// TODO: have a failsafe here?
		for i, r := range routes {
			minIndex := -1
			var min float64 = math.MaxFloat64
			for k, l := range remainingLoads {
				lCost := r.CompletionCostWithLoad(l)
				if lCost < min && lCost <= 12*60 {
					min = math.Min(min, lCost)
					minIndex = k
				}
			}

			if min < loadCost {
				driver = i
				load = minIndex
				loadCost = min
			}
		}

		if internal.Debug {
			fmt.Println(fmt.Sprintf("iter %v (%v): driver %v with load %v for cost %v", iter, len(remainingLoads), driver, load, loadCost))
		}

		if iter > 10 {
			break
		}
		routes[driver] = append(routes[driver], remainingLoads[load])
		delete(remainingLoads, load)
	}

	return routes
}
