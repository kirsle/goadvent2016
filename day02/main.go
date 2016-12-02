package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Line represents a line of steps from the input file.
type Line struct {
	Moves []Direction
}

// Direction represents a cardinal direction (Up, Right, Down, Left)
type Direction int

// Direction constants.
const (
	Up Direction = iota
	Right
	Down
	Left
)

// The keypad.
var Keypad = [5][5]string{
	{" ", " ", "1", " ", " "},
	{" ", "2", "3", "4", " "},
	{"5", "6", "7", "8", "9"},
	{" ", "A", "B", "C", " "},
	{" ", " ", "D", " ", " "},
}

// The dimensions of the keypad.
const KeypadSize = 5

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	// Get the list of steps to follow.
	lines, err := ParseInput(os.Args[1])
	if err != nil {
		panic(err)
	}

	// For our keypad map, (1,1) will represent the number '5' in the middle of
	// the pad, and the boundaries are therefore (0,0) to (2,2). We'll work
	// with these numbers and then use a converter function to tell us what
	// digit is at a given coordinate.
	x, y, passcode := 0, 2, make([]string, len(lines))
	_ = passcode

	Debug("Start at number: %s (at %d,%d)\n", GetNumber(x, y), x, y)

	// Check each line of instructions.
	for i, line := range lines {
		for _, move := range line.Moves {
			// Move our pointer.
			success := MovePointer(&x, &y, move)

			Debug("Line %d: move %d to position (%d,%d) - valid: %v - on key: %s\n", i, move, x, y, success, GetNumber(x, y))
		}

		// What digit is here?
		digit := GetNumber(x, y)
		Debug("Got pass code digit: %s\n", digit)
		passcode[i] = digit
	}

	fmt.Printf("The pass code is: %v\n", passcode)
}

// MovePointer attempts to move the pointer by 1 in a given direction, with
// bounds checking so it won't move into an invalid space. Returns true if the
// move was acceptable.
func MovePointer(x, y *int, d Direction) bool {
	// Check each direction of movement, see whether it's on the board and
	// there's a valid digit there, and move our coordinates if its OK.
	if d == Up && *y-1 >= 0 && CanMove(*x, *y-1) {
		*y -= 1
		return true
	} else if d == Right && *x+1 < KeypadSize && CanMove(*x+1, *y) {
		*x += 1
		return true
	} else if d == Down && *y+1 < KeypadSize && CanMove(*x, *y+1) {
		*y += 1
		return true
	} else if d == Left && *x-1 >= 0 && CanMove(*x-1, *y) {
		*x -= 1
		return true
	}

	return false
}

// GetNumber returns the number at the given coordinate.
func GetNumber(x, y int) string {
	return Keypad[y][x]
}

// CanMove returns whether the coordinate is an actual number on the keypad.
func CanMove(x, y int) bool {
	return Keypad[y][x] != " "
}

// ParseInput parses the input text file and returns an array of Steps.
func ParseInput(file string) ([]Line, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	// Make a line scanner.
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	// Make the buffer of lines.
	lines := []Line{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// Look for the individual direction steps on this line.
		row := Line{}
		for i := 0; i < len(line); i++ {
			char := line[i]
			if char == 'U' {
				row.Moves = append(row.Moves, Up)
			} else if char == 'R' {
				row.Moves = append(row.Moves, Right)
			} else if char == 'D' {
				row.Moves = append(row.Moves, Down)
			} else if char == 'L' {
				row.Moves = append(row.Moves, Left)
			} else {
				return nil, errors.New(fmt.Sprintf("Unexpected character in input file: %v", char))
			}
		}

		lines = append(lines, row)
	}

	return lines, nil
}

// Debug prints a debug message.
func Debug(template string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(template, a...)
	}
}
