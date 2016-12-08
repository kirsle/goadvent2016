package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kirsle/goadvent2016/advent"
)

const (
	ScreenWidth  = 50
	ScreenHeight = 6
)

// Type Screen represents our 50x6 LCD screen. The pixels are represented as an
// array of rows (Y coord) with each row having the X coord.
type Screen struct {
	Pixels [ScreenHeight][ScreenWidth]bool
}

// Regexps for the screen instructions.
var (
	RectangleRegexp *regexp.Regexp = regexp.MustCompile(`^rect (\d+)x(\d+)$`)
	ColumnRegexp    *regexp.Regexp = regexp.MustCompile(`^rotate column x=(\d+) by (\d+)$`)
	RowRegexp       *regexp.Regexp = regexp.MustCompile(`^rotate row y=(\d+) by (\d+)$`)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	input, err := advent.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Create our screen.
	screen := NewScreen()

	// Process the instructions.
	for _, instruction := range input {
		err = screen.ProcessInstruction(instruction)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}

		// When debugging, print the screen after every update.
		if os.Getenv("DEBUG") != "" {
			screen.Print()
		}
	}

	// Print the final screen.
	fmt.Println("Final screen:")
	screen.Print()

	// Count the lit pixels.
	fmt.Printf("Number of pixels lit: %d\n", screen.LitCount())
}

// NewScreen creates a new LCD screen with all the pixels turned off.
func NewScreen() *Screen {
	return &Screen{
		Pixels: [ScreenHeight][ScreenWidth]bool{},
	}
}

// Light turns on a pixel at a given coordinate.
func (s *Screen) Light(x, y int) error {
	err := s.BoundsCheck(x, y)
	if err == nil {
		s.Pixels[y][x] = true
	}
	return err
}

// Dark turns off a pixel at a given coordinate.
func (s *Screen) Dark(x, y int) error {
	err := s.BoundsCheck(x, y)
	if err == nil {
		s.Pixels[y][x] = false
	}
	return err
}

// IsLit tells whether a pixel at a given coordinate is lit.
func (s *Screen) IsLit(x, y int) (bool, error) {
	if err := s.BoundsCheck(x, y); err != nil {
		return false, err
	}
	return s.Pixels[y][x], nil
}

// LitCount counts the number of lit pixels.
func (s *Screen) LitCount() int {
	var count int
	for _, y := range s.Pixels {
		for _, x := range y {
			if x {
				count++
			}
		}
	}
	return count
}

// BoundsCheck checks whether an X and Y coordinate is valid.
func (s *Screen) BoundsCheck(x, y int) error {
	if y < 0 || y > ScreenHeight {
		return errors.New("Can't darken pixel: Y coordinate is out of bounds")
	} else if x < 0 || x > ScreenWidth {
		return errors.New("Can't darken pixel: X coordinate is out of bounds")
	} else {
		return nil
	}
}

// Print shows what the screen looks like in ASCII art.
func (s *Screen) Print() {
	for _, row := range s.Pixels {
		for _, col := range row {
			if col {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

// ProcessInstruction parses and executes a pixel manipulation function.
func (s *Screen) ProcessInstruction(input string) error {
	advent.Debug("INSTRUCTION: %s\n", input)

	// Check the type of instruction we're dealing with.
	if strings.HasPrefix(input, "rect") {
		return s.ProcessRect(input)
	} else if strings.HasPrefix(input, "rotate column") {
		return s.ProcessColumn(input)
	} else if strings.HasPrefix(input, "rotate row") {
		return s.ProcessRow(input)
	}
	return fmt.Errorf("Invalid instruction: %s", input)
}

// ProcessRect handles a rectangle drawing instruction.
func (s *Screen) ProcessRect(input string) error {
	match := RectangleRegexp.FindStringSubmatch(input)
	if len(match) == 0 {
		return fmt.Errorf("Failed regexp for rectangle: %s", input)
	}

	// Turn the regexp matches into ints.
	values, err := advent.StringsToInts(match[1:])
	if err != nil {
		return err
	}
	width, height := values[0], values[1]

	// Fill in the pixels.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			err := s.Light(x, y)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ProcessColumn handles a column shift instruction.
func (s *Screen) ProcessColumn(input string) error {
	match := ColumnRegexp.FindStringSubmatch(input)
	if len(match) == 0 {
		return fmt.Errorf("Failed regexp for column: %s", input)
	}

	values, err := advent.StringsToInts(match[1:])
	if err != nil {
		return err
	}
	column, length := values[0], uint(values[1])

	// Get all the pixels of this column into a convenient array.
	var pixels []bool
	for y := 0; y < ScreenHeight; y++ {
		var lit bool
		lit, err = s.IsLit(column, y)
		if err != nil {
			return err
		}
		pixels = append(pixels, lit)
	}

	// Shift the pixels.
	pixels = ShiftRight(pixels, length)

	// Update their lit values.
	for y, lit := range pixels {
		if lit {
			err = s.Light(column, y)
		} else {
			err = s.Dark(column, y)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// ProcessRow handles a row shift instruction.
func (s *Screen) ProcessRow(input string) error {
	match := RowRegexp.FindStringSubmatch(input)
	if len(match) == 0 {
		return fmt.Errorf("Failed regexp row row: %s", input)
	}

	values, err := advent.StringsToInts(match[1:])
	if err != nil {
		return err
	}
	row, length := values[0], uint(values[1])

	// Get all the pixels of this row into a convenient array.
	var pixels []bool
	for x := 0; x < ScreenWidth; x++ {
		var lit bool
		lit, err = s.IsLit(x, row)
		if err != nil {
			return err
		}
		pixels = append(pixels, lit)
	}

	// Shift the pixels.
	pixels = ShiftRight(pixels, length)

	// Update their lit pixels.
	for x, lit := range pixels {
		if lit {
			err = s.Light(x, row)
		} else {
			err = s.Dark(x, row)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// ShiftRight shifts an array of booleans forward with wrap-around.
func ShiftRight(array []bool, steps uint) []bool {
	var last bool
	var i uint
	for i = 0; i < steps; i++ {
		last = array[len(array)-1]
		for j := len(array) - 1; j > 0; j-- {
			array[j] = array[j-1]
		}
		array[0] = last
	}
	return array
}
