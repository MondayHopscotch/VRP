package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var coordMatcher = regexp.MustCompile(`\((-?\d+\.?\d*),(-?\d+\.?\d*)\)`)

func parseLine(input string) (Load, error) {
	fields := strings.Fields(input)
	if len(fields) != 3 {
		return Load{}, fmt.Errorf("malformed load: expected 3 fields, received %v", len(fields))
	}

	loadNum, err := strconv.Atoi(fields[0])
	if err != nil {
		return Load{}, errors.Wrap(err, "load number unable to be parsed")
	}

	pickup, err := parsePoint(fields[1])
	if err != nil {
		return Load{}, errors.Wrap(err, "pickup point unable to be parsed")
	}

	dropoff, err := parsePoint(fields[2])
	if err != nil {
		return Load{}, errors.Wrap(err, "dropoff point unable to be parsed")
	}

	return Load{
		Number:  loadNum,
		Pickup:  pickup,
		Dropoff: dropoff,
	}, nil
}

func parsePoint(input string) (Point, error) {
	if !coordMatcher.MatchString(input) {
		return Point{}, fmt.Errorf("point does not match the regex %v", coordMatcher.String())
	}

	groups := coordMatcher.FindStringSubmatch(input)

	x, err := strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return Point{}, errors.Wrap(err, "point x coordinate unable to be parsed")
	}

	y, err := strconv.ParseFloat(groups[2], 64)
	if err != nil {
		return Point{}, errors.Wrap(err, "point y coordinate unable to be parsed")
	}

	return Point{
		X: x,
		Y: y,
	}, nil
}
