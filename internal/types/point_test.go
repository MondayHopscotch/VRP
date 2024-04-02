package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistance(t *testing.T) {
	p1 := Point{X: 5, Y: 10}
	p2 := Point{X: -5, Y: 35}

	// Just a sanity check on our math, accounting for potential floating point rounding issues with a small threshold
	assert.LessOrEqual(t, math.Abs(math.Sqrt(math.Pow(10, 2)+math.Pow(25, 2))-p1.DistanceTo(p2)), 0.000001)
}
