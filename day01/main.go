package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// For debug output: export DEBUG=1

// Type Step represents a step from the input file, which contains a direction
// (left or right) and the number of blocks to travel.
type Step struct {
	Direction Direction
	Steps     int
}

// Type Direction represents either Left or Right.
type Direction int

// Direction constants.
const (
	Left Direction = iota
	Right
)

// Type Facing represents what direction we're facing.
type Facing int

// Facing constants.
const (
	North Facing = iota
	East
	South
	West
)

// Type Coordinate logs an X and Y pair for places we've been.
type Coordinate struct {
	X int
	Y int
}

// Type Visited stores a map of coordinates we've been to.
type Visited map[Coordinate]bool

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	// Get the list of steps to follow.
	steps, err := ParseInput(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Track our offsets starting at 0,0 and facing north to find where the
	// directions lead to.
	x, y, facing := 0, 0, North

	// Also keep track of where we've been. The Easter Bunny HQ is at the first
	// location that we visit twice.
	visited := &Visited{}
	visited.Visit(Coordinate{0, 0})

	// Follow the steps.
	for _, step := range steps {
		// Turn first.
		facing.Turn(step.Direction)

		// Print output for debugging.
		Debug("Step: %v - Now Facing: %v - Coords: (%d,%d)\n", step, facing, x, y)

		done := visited.Travel(&x, &y, facing, step.Steps)
		if done {
			break
		}
	}

	// And our verdict is...
	distance := int(math.Abs(float64(x)) + math.Abs(float64(y)))
	fmt.Printf("The Easter Bunny HQ is %d blocks away.\n", distance)
}

// Turn calculates what direction we're facing.
func (f *Facing) Turn(direction Direction) {
	// Rotate our direction of facing first.
	if direction == Right {
		*f += 1
	} else {
		*f -= 1
	}

	// And bounds check it.
	if *f < North {
		*f = West
	} else if *f > West {
		*f = North
	}
}

// Travel moves our position along a vector and returns true if we've stepped
// over the same position twice.
func (v *Visited) Travel(x, y *int, facing Facing, distance int) bool {
	// Loop for the distance desired.
	for i := 0; i < distance; i++ {
		// Move our position along the vector.
		if facing == North {
			*y += 1
		} else if facing == East {
			*x += 1
		} else if facing == South {
			*y -= 1
		} else if facing == West {
			*x -= 1
		}

		// The coordinate we're currently looking at.
		coord := Coordinate{*x, *y}

		// Mark it as visited. This also tells us whether we stepped over the
		// spot twice, so we can return true if so.
		if v.Visit(coord) {
			return true
		}
	}

	return false
}

// Visit marks a spot we've visited and returns true if it's a duplicate spot.
func (v *Visited) Visit(c Coordinate) bool {
	Debug("Visit coord: %v\n", c)
	if _, ok := (*v)[c]; ok {
		fmt.Printf("We stepped back over our tracks at %v!\n", c)
		return true
	}
	(*v)[c] = true
	return false
}

// ParseInput parses the input text file and returns an array of Steps.
func ParseInput(file string) ([]Step, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	// Make a line scanner.
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	// Make the buffer of steps.
	steps := []Step{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// Look for steps. Steps look like "R5" or "L2": a direction and a
		// number of blocks to travel that direction.
		for _, step := range strings.Split(line, ",") {
			step = strings.TrimSpace(step)
			direction := step[0]
			blocks, err := strconv.Atoi(step[1:])
			if err != nil {
				return nil, err
			}

			if direction == 'R' {
				steps = append(steps, Step{Right, blocks})
			} else if direction == 'L' {
				steps = append(steps, Step{Left, blocks})
			} else {
				return nil, errors.New(fmt.Sprintf("Found an invalid step entry: %v", step))
			}
		}
	}

	return steps, nil
}

// Debug prints a debug message.
func Debug(template string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(template, a...)
	}
}
