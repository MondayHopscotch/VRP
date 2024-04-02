package routing

import (
	"fmt"
	"math"
	"vrp/internal"
	"vrp/internal/types"
)

func NewBruteForceSolver(loads []types.Load) Solver {
	return &BruteForceSolver{loads: loads}
}

// BruteForceSolver plans routes basedon a simple greedy algorithm for routing our drives
type BruteForceSolver struct {
	loads []types.Load
}

func (b *BruteForceSolver) PlanRoutes() []types.Route {
	// TODO: require at least one driver
	// TODO: Need to return error if we can't find a set of routes

	// Add all drivers to our pool
	// Loop:
	//  Find closest load to each driver that *doesn't put them into overtime with the trip to hub taken into account*
	//      Note that this may not be the closest load, as something further away, but also on the way to the hub may be a doable option
	//  Assign highest priority load to the appropriate driver
	//
	// We have to be careful not to strand drivers

	// TODO: determine number of drivers required
	routes := make([]types.Route, 4)

	remainingLoads := make(map[int]types.Load)
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
				lCost := r.CompletionTimeWithLoad(l)
				if lCost < min && lCost <= types.DriverMaxTime {
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
