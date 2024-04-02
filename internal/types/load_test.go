package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCost(t *testing.T) {
	load := Load{
		Number:  1,
		Pickup:  Point{X: 10, Y: 10},
		Dropoff: Point{X: 50, Y: -10},
	}

	assert.Equal(t, math.Hypot(40, 20), load.Cost())
}
