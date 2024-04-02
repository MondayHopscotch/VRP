package routing

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"vrp/internal"
)

type NearestNeighborSolver struct {
	loads []internal.Load
}

func NewNearestNeighborSolver(loads []internal.Load) Solver {
	solver := NearestNeighborSolver{
		loads: loads,
	}

	// Add our depot as a special load that starts and ends at the depot
	solver.loads = append([]internal.Load{internal.Load{Number: 0, Pickup: internal.Point{}, Dropoff: internal.Point{}}}, solver.loads...)

	if internal.Debug {
		fmt.Println(fmt.Sprintf("Nearest Neighbor Solver built with %v loads", len(loads)))
	}
	return solver
}

func (n NearestNeighborSolver) PlanRoutes() []internal.Route {
	neighbors := n.getNeighborMap()

	if internal.Debug {
		for i := 0; i < len(neighbors); i++ {
			for _, o := range neighbors[i] {
				fmt.Println(fmt.Sprintf("\t(%v)\t%v", n.loads[i].Dropoff.DistanceTo(o.Pickup), o))
			}
		}
	}

	roughMinTotal := 0.0
	visited := []int{0}
	currentLoadIndex := 0
	var current internal.Load
	for len(visited) < len(n.loads) {
		current = n.loads[currentLoadIndex]
		for _, l := range neighbors[currentLoadIndex] {
			if !slices.Contains(visited, l.Number) {
				roughMinTotal += current.Dropoff.DistanceTo(l.Pickup)
				currentLoadIndex = l.Number
				visited = append(visited, currentLoadIndex)
				break
			}
		}
	}

	// Add our final return to depot
	roughMinTotal += current.Dropoff.DistanceTo(n.loads[0].Pickup)

	// Each driver can only drive for 12 hours
	minDrivers := math.Ceil(roughMinTotal / (12 * 60))

	if internal.Debug {
		fmt.Println(visited)
		fmt.Println(roughMinTotal)
		fmt.Println(minDrivers)
	}

	routes := make([]internal.Route, int(minDrivers))
	for i, _ := range routes {
		routes[i] = append(routes[i], internal.Load{})
	}

	remainingLoads := make(map[int]internal.Load)
	for _, l := range n.loads {
		if l.Number == 0 {
			// we don't need to track our hub
			continue
		}
		remainingLoads[l.Number] = l
	}

	var driverFound bool
	var driver int
	var nextLoad internal.Load
	var nearestCost float64
	for len(remainingLoads) > 0 {
		driverFound = false
		nearestCost = math.MaxFloat64
		for i, r := range routes {
			for _, l := range neighbors[r[len(r)-1].Number] {
				if _, ok := remainingLoads[l.Number]; ok {
					costIncrease := r.CostIncreaseWithLoad(l)
					// Check driver capacity
					newRouteCostCheck := r.CompletionCostWithLoad(l)
					if costIncrease < nearestCost && newRouteCostCheck <= 12*60 {
						driverFound = true
						driver = i
						nextLoad = l
						nearestCost = costIncrease

						// As our neighbors are already sorted by nearest, we minimize deadhead by
						// going with the first match here and moving on
						break
					}
				}
			}
		}

		if driverFound {
			routes[driver] = append(routes[driver], nextLoad)
			delete(remainingLoads, nextLoad.Number)
		} else {
			if internal.Debug {
				fmt.Println(fmt.Sprintf("no driver found for remaining loads (%v). Adding driver", remainingLoads))
			}
			routes = append(routes, internal.Route{internal.Load{}})
		}
	}

	// Prune out our depot 'loads'
	for i, r := range routes {
		routes[i] = slices.Delete(r, 0, 1)
	}

	return routes
}

func (n NearestNeighborSolver) getNeighborMap() map[int][]internal.Load {
	neighbors := make(map[int][]internal.Load)

	for _, l := range n.loads {
		closest := make([]internal.Load, 0)
		for _, o := range n.loads {
			if o == l {
				continue
			}
			closest = append(closest, o)
		}
		sort.Slice(closest, func(a, b int) bool {
			nextA := closest[a]
			nextB := closest[b]
			return l.Dropoff.DistanceTo(nextA.Pickup) < l.Dropoff.DistanceTo(nextB.Pickup)
		})
		neighbors[l.Number] = closest
	}

	return neighbors
}
