package routing

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"vrp/internal"
	"vrp/internal/types"
)

type NearestNeighborSolver struct {
	loads []types.Load
}

func NewNearestNeighborSolver(loads []types.Load) Solver {
	solver := NearestNeighborSolver{
		loads: loads,
	}

	// Add our depot as a special load that starts and ends at the depot
	solver.loads = append([]types.Load{types.Load{Number: 0, Pickup: types.HubPoint(), Dropoff: types.Point{}}}, solver.loads...)

	if internal.Debug {
		fmt.Println(fmt.Sprintf("Nearest Neighbor Solver built with %v loads", len(loads)))
	}
	return solver
}

func (n NearestNeighborSolver) PlanRoutes() []types.Route {
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
	var current types.Load
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

	// Each driver has max shift length
	minDrivers := int(math.Ceil(roughMinTotal / types.DriverMaxTime))

	if internal.Debug {
		fmt.Println(visited)
		fmt.Println(roughMinTotal)
		fmt.Println(minDrivers)
	}

	resultRoutes, totalCost := n.planRoutesForDrivers(minDrivers, neighbors)

	if len(resultRoutes) != minDrivers {
		if internal.Debug {
			fmt.Println(fmt.Sprintf("drivers were added to complete loads (%v -> %v)", minDrivers, len(resultRoutes)))
		}

		// we had to add drivers. Recalculate and see if starting with extra drivers yields cost improvement
		newRoutes, newTotalCost := n.planRoutesForDrivers(len(resultRoutes), neighbors)
		if internal.Debug {
			fmt.Println(fmt.Sprintf("recalculation with %v drivers yielded cost %v (previous %v", len(resultRoutes), newTotalCost, totalCost))
		}
		if newTotalCost < totalCost {
			resultRoutes = newRoutes
		}
	}

	// Prune out our depot 'loads'
	for i, r := range resultRoutes {
		resultRoutes[i] = slices.Delete(r, 0, 1)
	}

	return resultRoutes
}

func (n NearestNeighborSolver) planRoutesForDrivers(startingDrivers int, neighbors map[int][]types.Load) ([]types.Route, float64) {
	routes := make([]types.Route, startingDrivers)
	for i, _ := range routes {
		routes[i] = append(routes[i], types.Load{})
	}

	remainingLoads := make(map[int]types.Load)
	for _, l := range n.loads {
		if l.Number == 0 {
			// we don't need to track our hub
			continue
		}
		remainingLoads[l.Number] = l
	}

	var driverFound bool
	var driver int
	var nextLoad types.Load
	for len(remainingLoads) > 0 {
		driverFound = false
		for i, r := range routes {
			for _, l := range neighbors[r[len(r)-1].Number] {
				if _, ok := remainingLoads[l.Number]; ok {
					// Check driver capacity
					if r.CompletionTimeWithLoad(l) <= types.DriverMaxTime {
						driverFound = true
						driver = i
						nextLoad = l

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
				fmt.Println(fmt.Sprintf("no driver found for remaining loads (%v). Adding driver", len(remainingLoads)))
			}
			routes = append(routes, types.Route{types.Load{}})
		}
	}

	totalCost := 0.0

	for _, r := range routes {
		totalCost += r.TotalCostWithDriver()
	}

	return routes, totalCost
}

func (n NearestNeighborSolver) getNeighborMap() map[int][]types.Load {
	neighbors := make(map[int][]types.Load)

	for _, l := range n.loads {
		closest := make([]types.Load, 0)
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
