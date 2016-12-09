package main

import (
	"strings"
	"testing"
)

// Type TestCase contains test cases.
type TestCase struct {
	Input          string // The input to the algorithm.
	ExpectedOutput string // The expected result (don't check if empty)
	ExpectedLength int    // The expected length of the result
	ShouldError    bool   // Whether we're expecting an error or not
}

func init() {
	// We want the compressed data returned to verify the algorithms work.
	ReturnData = true
}

func TestDecompressV1(t *testing.T) {
	tests := []TestCase{
		TestCase{"ADVENT", "ADVENT", 6, false},
		TestCase{"A(1x5)BC", "ABBBBBC", 7, false},
		TestCase{"(3x3)XYZ", "XYZXYZXYZ", 9, false},
		TestCase{"A(2x2)BCD(2x2)EFG", "ABCBCDEFEFG", 11, false},
		TestCase{"(6x1)(1x3)A", "(1x3)A", 6, false},
		TestCase{"X(8x2)(3x3)ABCY", "X(3x3)ABC(3x3)ABCY", 18, false},
	}
	runTestCases(t, tests, 1)
}

func TestDecompressV2(t *testing.T) {
	tests := []TestCase{
		// Good test cases.
		TestCase{"ADVENT", "ADVENT", 6, false},
		TestCase{"A(1x5)BC", "ABBBBBC", 7, false},
		TestCase{"(3x3)XYZ", "XYZXYZXYZ", 9, false},
		TestCase{"X(8x2)(3x3)ABCY", "XABCABCABCABCABCABCY", 20, false},
		TestCase{"(27x12)(20x12)(13x14)(7x10)(1x12)A", strings.Repeat("A", 241920), 241920, false},
		TestCase{"(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN", "", 445, false},
	}
	runTestCases(t, tests, 2)
}

func runTestCases(t *testing.T, tests []TestCase, version int) {
	for _, test := range tests {
		output, size, err := Decompress(test.Input, version)
		if err != nil {
			if !test.ShouldError {
				t.Errorf("Unexpected error from test: %v", err)
			}
			continue
		}

		if len(test.ExpectedOutput) > 0 && output != test.ExpectedOutput {
			t.Errorf(`Output assertion error: expected "%s", got "%s"`, test.ExpectedOutput, output)
		}
		if size != test.ExpectedLength {
			t.Errorf(`Output length assertion error: expected "%d", got "%d"`, test.ExpectedLength, size)
		}
	}
}
