package parsing

import (
	"strings"
	"testing"
	"vrp/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestParsePoint(t *testing.T) {
	point, err := parsePoint("(123.456,456.789)")
	assert.Nil(t, err)
	AssertPoint(t, types.Point{X: 123.456, Y: 456.789}, point)
}

func TestParseBadPoint(t *testing.T) {
	_, err := parsePoint("( ,456.789)")
	assert.ErrorContains(t, err, "point does not match the regex")
}

func TestSimpleLineParse(t *testing.T) {
	line := "1 (-50.1,80.0) (90.1,12.2)"
	load, err := parseLine(line)
	assert.Nil(t, err)
	assert.Equal(t, 1, load.Number)
	AssertPoint(t, load.Pickup, types.Point{
		X: -50.1,
		Y: 80.0,
	})
	AssertPoint(t, load.Dropoff, types.Point{
		X: 90.1,
		Y: 12.2,
	})
}

func TestMalformattedLineParse(t *testing.T) {
	line := "(1.1,2.2) (5,10.5)"
	_, err := parseLine(line)
	assert.ErrorContains(t, err, "malformed load")
}

func TestBadLoadNumber(t *testing.T) {
	_, err := parseLine(`hi (1,1) (2,2)`)
	assert.ErrorContains(t, err, "load number unable to be parsed")
}

func TestBadPickupCoordinate(t *testing.T) {
	_, err := parseLine(`5 (j) (2,2)`)
	assert.ErrorContains(t, err, "pickup point unable to be parsed")
}

func TestBadDropoffCoordinate(t *testing.T) {
	_, err := parseLine(`5 (1,1) (z)`)
	assert.ErrorContains(t, err, "dropoff point unable to be parsed")
}

func TestParseReader(t *testing.T) {
	contents := `loadNumber pickup dropoff
    1 (1.25,2.50) (3.75,4.0)`
	loads, err := parseReader(strings.NewReader(contents))
	assert.Nil(t, err)
	assert.Len(t, loads, 1)
	AssertPoint(t, types.Point{X: 1.25, Y: 2.5}, loads[0].Pickup)
	AssertPoint(t, types.Point{X: 3.75, Y: 4.0}, loads[0].Dropoff)
}

// AssertPoint asserts that the actual Point's X and Y match the expected values
func AssertPoint(t *testing.T, expected types.Point, actual types.Point) {
	assert.Equal(t, expected.X, actual.X)
	assert.Equal(t, expected.Y, actual.Y)
}
