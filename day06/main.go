package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Type Position keeps track of most frequently seen letters in a given position.
type Position struct {
	Frequency map[rune]int
}

// Type Code keeps track of various letters in various positions to determine
// the repetition code's message.
type Code struct {
	Positions []*Position
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	// The eventual code we're trying to crack.
	code := &Code{}

	// Get the input strings.
	inputs := ReadInputFile(os.Args[1])
	for _, input := range inputs {
		// Initialize the length of the code.
		for len(code.Positions) < len(input) {
			code.Positions = append(code.Positions, NewPosition())
		}

		// Check each position of the input.
		for i, char := range input {
			code.Positions[i].Put(char)
		}
	}

	// Get the most/least frequent symbols (Part 1 & 2 of the puzzle)
	fmt.Printf("The most likely code is: %s\n", code.MostLikely())
	fmt.Printf("The least likely is: %s\n", code.LeastLikely())
}

// NewPosition initializes a new position object.
func NewPosition() *Position {
	new := &Position{
		Frequency: map[rune]int{},
	}
	return new
}

// Put adds a letter to the position and increments the letter count.
func (p *Position) Put(letter rune) {
	_, ok := p.Frequency[letter]
	if !ok {
		p.Frequency[letter] = 0
	}
	p.Frequency[letter]++
}

// Most returns the most frequently used letter in a position.
func (p *Position) Most() rune {
	var result rune
	var highest int

	for k, v := range p.Frequency {
		if v > highest {
			result = k
			highest = v
		}
	}

	return result
}

// Least returns the least frequently used letter in a position.
func (p *Position) Least() rune {
	var result rune
	var lowest int

	for k, v := range p.Frequency {
		if lowest == 0 || v < lowest {
			result = k
			lowest = v
		}
	}

	return result
}

// MostLikely shows the string value of the code's most likely symbols.
func (c *Code) MostLikely() string {
	result := []string{}
	for _, p := range c.Positions {
		result = append(result, string(p.Most()))
	}
	return strings.Join(result, "")
}

// LeastLikely shows the string value of the code's least likely symbols.
func (c *Code) LeastLikely() string {
	result := []string{}
	for _, p := range c.Positions {
		result = append(result, string(p.Least()))
	}
	return strings.Join(result, "")
}

// ReadInputFile produces the list of strings from the file.
func ReadInputFile(filename string) []string {
	fh, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	// Lines read.
	parsed := []string{}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		parsed = append(parsed, line)
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}

	return parsed
}
