package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var WhitespaceRegexp = regexp.MustCompile(`\s+`)

// Type Line represents a literal line of numbers from the input file.
type Line struct {
	A int
	B int
	C int
}

// Type Triangle represents three dimensions of a triangle's edges.
type Triangle struct {
	A int
	B int
	C int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	// Parse the input file into an array of lines of numbers.
	inputLines := ParseInput(os.Args[1])

	// But the triangles in the input are arranged vertically! If this weren't
	// the case we'd simply use `inputLines` as the `triangles` list.
	triangles := ParseTriangles(inputLines)

	// Look for invalid triangles.
	var invalid int = 0
	for i, triangle := range triangles {
		log.Printf("Triangle: %v", triangle)
		if !triangle.IsValid() {
			invalid++
			log.Printf("Invalid triangle at line %d: %v", i, triangle)
		}
	}

	fmt.Printf("Of %d triangles, %d are valid and %d are not valid\n",
		len(triangles),
		len(triangles)-invalid,
		invalid,
	)
}

// ParseTriangles turns the input lines of numbers into triangles.
func ParseTriangles(lines []Line) []Triangle {
	result := []Triangle{}

	// We need to scan the input 3 lines at a time and produce 3 triangles
	// from each column.
	for i := 0; i < len(lines); i += 3 {
		rows := lines[i : i+3]

		result = append(result,
			Triangle{rows[0].A, rows[1].A, rows[2].A},
			Triangle{rows[0].B, rows[1].B, rows[2].B},
			Triangle{rows[0].C, rows[1].C, rows[2].C},
		)
	}

	return result
}

// IsValid validates a triangle.
func (t Triangle) IsValid() bool {
	// Test all the permutations of sides to see if the triangle is invalid.
	permutations := []bool{
		t.A+t.B > t.C,
		t.A+t.C > t.B,
		t.B+t.C > t.A,
	}
	for _, ok := range permutations {
		if !ok {
			return false
		}
	}

	return true
}

// ParseInput parses the lines of integers from the input file.
func ParseInput(file string) []Line {
	fh, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	// The triangles parsed from the file.
	result := []Line{}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		// Convert the numbers on this line to ints.
		numbers := WhitespaceRegexp.Split(strings.TrimSpace(scanner.Text()), 3)
		sides := []int{}
		for _, side := range numbers {
			value, err := strconv.Atoi(side)
			if err != nil {
				log.Println(err)
				continue
			}
			sides = append(sides, value)
		}

		// Make sure we got a valid triplet
		if len(sides) != 3 {
			log.Println("Didn't get a valid line set")
			continue
		}

		result = append(result, Line{sides[0], sides[1], sides[2]})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}
