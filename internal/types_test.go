package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteCost(t *testing.T) {
	route := Route{}
	route = append(route, Load{
		Number:  1,
		Pickup:  Point{X: 3, Y: 0},
		Dropoff: Point{X: 3, Y: 4},
	})

	// Using a simple 3 4 5 triangle because it's easy
	assert.Equal(t, 500.0+3+4, route.CurrentCost())
	assert.Equal(t, 500.0+3+4+5, route.CurrentCompletionCost())
}

func TestRouteCostWithLoad(t *testing.T) {
	route := Route{}

	load := Load{
		Number:  1,
		Pickup:  Point{X: 3, Y: 0},
		Dropoff: Point{X: 3, Y: 4},
	}

	// Using a simple 3 4 5 triangle because it's easy
	assert.Equal(t, 500.0+3+4, route.CostWithLoad(load))
	assert.Equal(t, 500.0+3+4+5, route.CompletionCostWithLoad(load))
}
