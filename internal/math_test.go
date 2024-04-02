package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTravelCost(t *testing.T) {
	start := Point{X: 0.0, Y: 0.0}
	end := Point{X: 2.0, Y: 3.0}

	assert.Equal(t, 4.0, TravelCost(start, end))
}
