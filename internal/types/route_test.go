package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteTime(t *testing.T) {
	route := Route{}
	assert.Equal(t, 0.0, route.CurrentTime())
	assert.Equal(t, 0.0, route.CurrentCompletionTime())

	route = append(route, Load{
		Number:  1,
		Pickup:  Point{X: 3, Y: 0},
		Dropoff: Point{X: 3, Y: 4},
	})

	// Using a simple 3 4 5 triangle because it's easy
	assert.Equal(t, 3.0+4, route.CurrentTime())
	assert.Equal(t, 3.0+4+5, route.CurrentCompletionTime())
}

func TestRouteCostWithDriver(t *testing.T) {
	route := Route{}
	assert.Equal(t, 500.0, route.TotalCostWithDriver())

	route = append(route, Load{
		Number:  1,
		Pickup:  Point{X: 3, Y: 0},
		Dropoff: Point{X: 3, Y: 4},
	})

	// Using a simple 3 4 5 triangle because it's easy
	assert.Equal(t, 500.0+3+4+5, route.TotalCostWithDriver())
}

func TestRouteWithLoad(t *testing.T) {
	route := Route{}

	load := Load{
		Number:  1,
		Pickup:  Point{X: 3, Y: 0},
		Dropoff: Point{X: 3, Y: 4},
	}

	assert.Equal(t, 7.0, route.TimeWithLoad(load))
	assert.Equal(t, 3.0+4+5, route.CompletionTimeWithLoad(load))
}

func TestRouteTimeIncrease(t *testing.T) {
	route := Route{}

	route = append(route, Load{
		Number:  1,
		Pickup:  Point{X: 3, Y: 0},
		Dropoff: Point{X: 3, Y: 3},
	})

	load := Load{
		Number:  1,
		Pickup:  Point{X: -3, Y: 3},
		Dropoff: Point{X: -3, Y: 0},
	}

	returnLegOne := math.Hypot(3, 3)
	totalTravelOneLoad := 6.0 + returnLegOne

	returnLegTwo := 3.0
	totalTravelTwoLoad := 6.0 + 9.0 + returnLegTwo

	increase := totalTravelTwoLoad - totalTravelOneLoad

	assert.Equal(t, totalTravelOneLoad, route.CurrentCompletionTime())
	assert.Equal(t, increase, route.CompletionTimeIncreaseWithLoad(load))
}

func TestEmptyRoutesCost(t *testing.T) {
	routes := []Route{}
	assert.Equal(t, 0.0, GetTotalCostOfRoutes(routes))
}

func TestSimpleRoutesCost(t *testing.T) {
	routes := []Route{
		[]Load{
			{
				Number:  0,
				Pickup:  HubPoint(),
				Dropoff: HubPoint(),
			},
			{
				Number:  6,
				Pickup:  Point{X: 2, Y: 0},
				Dropoff: Point{X: 2, Y: 7},
			},
			{
				Number:  1,
				Pickup:  Point{X: 1, Y: 3},
				Dropoff: Point{X: -1, Y: -8},
			},
		},
		[]Load{
			{
				Number:  0,
				Pickup:  HubPoint(),
				Dropoff: HubPoint(),
			},
			{
				Number:  5,
				Pickup:  Point{X: 2, Y: 2},
				Dropoff: Point{X: -4, Y: 1},
			},
		},
	}

	assert.Equal(t, routes[0].TotalCostWithDriver()+routes[1].TotalCostWithDriver(), GetTotalCostOfRoutes(routes))
}
