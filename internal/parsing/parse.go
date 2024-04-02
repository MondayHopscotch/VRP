package parsing

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"vrp/internal"
	"vrp/internal/types"

	"github.com/pkg/errors"
)

// pointMatcher is a regex to match parentheses-surrounded x,y coordinates of floating point numbers. Ex: `(1.5,2.9)`
var pointMatcher = regexp.MustCompile(`\((-?\d+\.?\d*),(-?\d+\.?\d*)\)`)

// ParseAllLoads reads the file at the provided path and parses out a slice of Loads
func ParseAllLoads(inFile string) ([]types.Load, error) {
	file, err := os.Open(inFile)
	if err != nil {
		return []types.Load{}, errors.Wrap(err, "failed to read load file")
	}
	defer file.Close()

	return parseReader(file)
}

func parseReader(reader io.Reader) ([]types.Load, error) {
	loads := make([]types.Load, 0)

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum == 1 {
			// throw out header line
			continue
		}

		if internal.Debug {
			fmt.Println(fmt.Sprintf("parsing line: %v", scanner.Text()))
		}

		l, err := parseLine(scanner.Text())
		if err != nil {
			return []types.Load{}, errors.Wrap(err, fmt.Sprintf("unable to parse line %v", lineNum))
		}

		loads = append(loads, l)
	}

	return loads, nil
}

func parseLine(input string) (types.Load, error) {
	fields := strings.Fields(input)
	if len(fields) != 3 {
		return types.Load{}, fmt.Errorf("malformed load: expected 3 fields, received %v", len(fields))
	}

	loadNum, err := strconv.Atoi(fields[0])
	if err != nil {
		return types.Load{}, errors.Wrap(err, "load number unable to be parsed")
	}

	pickup, err := parsePoint(fields[1])
	if err != nil {
		return types.Load{}, errors.Wrap(err, "pickup point unable to be parsed")
	}

	dropoff, err := parsePoint(fields[2])
	if err != nil {
		return types.Load{}, errors.Wrap(err, "dropoff point unable to be parsed")
	}

	return types.Load{
		Number:  loadNum,
		Pickup:  pickup,
		Dropoff: dropoff,
	}, nil
}

func parsePoint(input string) (types.Point, error) {
	if !pointMatcher.MatchString(input) {
		return types.Point{}, fmt.Errorf("point does not match the regex %v", pointMatcher.String())
	}

	groups := pointMatcher.FindStringSubmatch(input)

	x, err := strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return types.Point{}, errors.Wrap(err, "point x coordinate unable to be parsed")
	}

	y, err := strconv.ParseFloat(groups[2], 64)
	if err != nil {
		return types.Point{}, errors.Wrap(err, "point y coordinate unable to be parsed")
	}

	return types.Point{
		X: x,
		Y: y,
	}, nil
}
