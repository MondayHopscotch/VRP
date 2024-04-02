package routing

import (
	"testing"
	"vrp/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestConstruction(t *testing.T) {
	solver := NewNearestNeighborSolver([]types.Load{})
	assert.NotNil(t, solver)
}

func TestEmptyLoads(t *testing.T) {
	solver := NewNearestNeighborSolver([]types.Load{})

	routes, err := solver.PlanRoutes()
	assert.Nil(t, err)
	assert.Empty(t, routes)
}

func TestSingleLoad(t *testing.T) {
	loadList := []types.Load{
		{
			Number:  1,
			Pickup:  types.Point{X: 4, Y: 3},
			Dropoff: types.Point{X: -4, Y: 3},
		},
	}

	solver := NewNearestNeighborSolver(loadList)

	routes, err := solver.PlanRoutes()
	assert.Nil(t, err)
	assert.Len(t, routes, 1)
	assert.Equal(t, types.Route(loadList), routes[0])
}

func TestSingleLoadOutOfRange(t *testing.T) {
	loadList := []types.Load{
		{
			Number:  1,
			Pickup:  types.Point{X: 0, Y: types.DriverMaxTime/2 + 1},
			Dropoff: types.Point{X: 0, Y: 1},
		},
	}

	solver := NewNearestNeighborSolver(loadList)

	routes, err := solver.PlanRoutes()
	assert.Empty(t, routes)
	assert.ErrorContains(t, err, "too far away")
}

func TestNeighborMap(t *testing.T) {
	loads := []types.Load{
		{
			Number:  1,
			Pickup:  types.Point{X: 3, Y: 3},
			Dropoff: types.Point{X: 3, Y: 5},
		},
		{
			Number:  2,
			Pickup:  types.Point{X: 0, Y: 5},
			Dropoff: types.Point{X: -1, Y: -1},
		},
	}

	solver := NewNearestNeighborSolver(loads)

	nMap := solver.getNeighborMap()
	assert.Len(t, nMap, 3)
	assert.Contains(t, nMap, 0)
	assert.Len(t, nMap[0], 2)
	assert.Equal(t, loads[0], nMap[0][0])
	assert.Equal(t, loads[1], nMap[0][1])

	assert.Contains(t, nMap, 1)
	assert.Len(t, nMap[1], 2)
	assert.Equal(t, loads[1], nMap[1][0])
	assert.Equal(t, types.HubPoint(), nMap[1][1].Pickup)

	assert.Contains(t, nMap, 2)
	assert.Equal(t, types.HubPoint(), nMap[2][0].Pickup)
	assert.Equal(t, loads[0], nMap[2][1])
}

func TestEstimateMinimum(t *testing.T) {
	loads := []types.Load{
		{
			Number:  1,
			Pickup:  types.Point{X: 3, Y: 3},
			Dropoff: types.Point{X: 3, Y: 5},
		},
		{
			Number:  2,
			Pickup:  types.Point{X: 0, Y: 5},
			Dropoff: types.Point{X: -1, Y: -1},
		},
	}

	solver := NewNearestNeighborSolver(loads)
	nMap := solver.getNeighborMap()

	min := solver.estimateMinimumTime(nMap)
	expected := types.HubPoint().DistanceTo(loads[0].Pickup)
	expected += loads[0].Cost()
	expected += loads[0].Dropoff.DistanceTo(loads[1].Pickup)
	expected += loads[1].Cost()
	expected += loads[1].Dropoff.DistanceTo(types.HubPoint())

	assert.Equal(t, expected, min)
}
