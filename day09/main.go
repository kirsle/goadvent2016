package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/kirsle/goadvent2016/advent"
)

const (
	// Sanity check for deep recursion prevention.
	RecursionLimit = 100

	// Version of the algorithm that main.go should use by default.
	AlgorithmVersion = 2
)

// If you actually want the string output, set this to true. For the input
// text on the V2 algorithm, the actual data is so big you'll run out of
// memory (10GB!).
//
// The unit tests set this to true to validate the algorithm.
var ReturnData = false

var MarkerRegexp *regexp.Regexp = regexp.MustCompile(`(\w*)\((\d+?)x(\d+?)\)`)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	var (
		input   []byte
		decoded string
		size    int
		err     error
	)

	input, err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	decoded, size, err = Decompress(strings.TrimSpace(string(input)), AlgorithmVersion)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded output: %s\nLength: %d\n", advent.Truncate(decoded, 255), size)
}

// Decompress implements the decompression algorithm.
func Decompress(input string, version int) (string, int, error) {
	advent.Debug("### INPUT: %s ###\n", input)

	// result is the actual string output if ReturnData is true.
	// totalSize is the size of the output whether or not we're actually
	// collecting it.
	result := bytes.NewBuffer([]byte{})
	var totalSize int

	// Put something in the result if we're not collecting one, for
	// debugging purposes.
	if !ReturnData {
		result.WriteString("<not returning data>")
	}

	// Index pointer into the input string that creeps along as we scan through.
	var idx int
	for idx < len(input) {
		advent.Debug("[%d] %s\n", idx, string(input[idx]))

		// Look for the next marker and catch any prefix characters before it.
		match := MarkerRegexp.FindStringSubmatch(input[idx:])

		// If no additional markers, glob up the remaining text and finish.
		if len(match) == 0 {
			advent.Debug("No more markers\n")
			if ReturnData {
				result.WriteString(input[idx:])
			}
			totalSize += len(input[idx:])
			break
		}

		advent.Debug("Found marker: %v\n", match)

		// Get the regexp parts separated.
		marker := match[0] // The full matched regexp including the prefix
		prefix := match[1] // Just the prefix part
		ints, _ := advent.StringsToInts(match[2:])
		length, repeat := ints[0], ints[1]

		// Shift the index past the marker.
		idx += len(marker)

		// The segment of text that needs repeating.
		var segment string
		var segSize int

		// Recursively expand.
		if version == 2 {
			advent.Debug("Descend recursively for: %s\n", input[idx:(idx+length)])
			subexpand, size, err := Decompress(input[idx:(idx+length)], version)
			if err != nil {
				return "", 0, err
			}
			segment = subexpand
			segSize = size
		} else {
			segment = input[idx : idx+length]
			segSize = len(segment)
		}

		// Glob in the prefix if any.
		if ReturnData {
			result.WriteString(prefix)
		}
		totalSize += len(prefix)

		// Run the repetition.
		advent.Debug("Repeat '%s' %d times\n", segment, repeat)
		for i := 0; i < repeat; i++ {
			if ReturnData {
				result.WriteString(segment)
			}
			totalSize += segSize
		}

		idx += length
	}

	advent.Debug("--- OUTPUT(len=%d): %s\n", totalSize, advent.Truncate(result.String(), 255))

	return result.String(), totalSize, nil
}
