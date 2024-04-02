package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePoint(t *testing.T) {
	point, err := parsePoint("(123.456,456.789)")
	assert.Nil(t, err)
	AssertPoint(t, Point{X: 123.456, Y: 456.789}, point)
}

func TestSimpleLineParse(t *testing.T) {
	line := "1 (-50.1,80.0) (90.1,12.2)"
	load, err := parseLine(line)
	assert.Nil(t, err)
	assert.Equal(t, 1, load.Number)
	AssertPoint(t, load.Pickup, Point{
		X: -50.1,
		Y: 80.0,
	})
	AssertPoint(t, load.Dropoff, Point{
		X: 90.1,
		Y: 12.2,
	})
}

func TestMalformattedLineParse(t *testing.T) {
	line := "(1.1,2.2) (5,10.5)"
	_, err := parseLine(line)
	assert.ErrorContains(t, err, "malformed load")
}

func AssertPoint(t *testing.T, expected Point, actual Point) {
	assert.Equal(t, expected.X, actual.X)
	assert.Equal(t, expected.Y, actual.Y)
}
