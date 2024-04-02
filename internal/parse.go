package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var coordMatcher = regexp.MustCompile(`\((-?\d+\.?\d*),(-?\d+\.?\d*)\)`)

func ParseAllLoads(inFile string) ([]Load, error) {
	file, err := os.Open(inFile)
	if err != nil {
		return []Load{}, errors.Wrap(err, "failed to read load file")
	}
	defer file.Close()

	loads := make([]Load, 0)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum == 1 {
			// throw out header line
			continue
		}

		if Debug {
			fmt.Println(fmt.Sprintf("parsing line: %v", scanner.Text()))
		}

		l, err := parseLine(scanner.Text())
		if err != nil {
			return []Load{}, errors.Wrap(err, fmt.Sprintf("unable to parse line %v", lineNum))
		}

		loads = append(loads, l)
	}

	return loads, nil
}

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
