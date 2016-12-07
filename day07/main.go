package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Type Address contains the parts of an IPv7 address.
type Address struct {
	Address  string   // The original string version of the address
	Supernet []string // The parts outside of square brackets
	Hypernet []string // The parts inside of square brackets
}

// BracketRegexp matches a string of characters in square brackets.
var BracketRegexp *regexp.Regexp = regexp.MustCompile(`\[([a-z]+?)\]`)

// SlidingWindow generalizes the procedure of scanning a sliding window for
// something you're looking for. It takes the input string, the size of the
// sliding window (e.g. 3 or 4 characters at a time), and a comparator function.
//
// The comparator is called for every slice in the input string. If the
// comparator returns true, this aborts the window and returns the slice that
// the comparator last looked at.
//
// See the implementation in HasAbba() and FindAbaBab()
func SlidingWindow(input string, size int, cmp func(string) bool) string {
	for i := 0; i+size-1 < len(input); i++ {
		slice := input[i : i+size]
		if cmp(slice) {
			return slice
		}
	}

	return ""
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	// Get the inputs.
	addresses := ParseAddresses(ReadFile(os.Args[1]))

	// Count the ones that support TLS and SSL.
	var supportsTLS int
	var supportsSSL int

	// Scan the addresses.
	for _, addr := range addresses {
		if addr.SupportsTLS() {
			Debug("%s supports TLS\n", addr.Address)
			supportsTLS++
		}
		if addr.SupportsSSL() {
			Debug("%s supports SSL\n", addr.Address)
			supportsSSL++
		}
	}

	fmt.Printf("%d addresses support TLS.\n", supportsTLS)
	fmt.Printf("%d addresses support SSL.\n", supportsSSL)
}

// NewAddress creates a new address object.
func NewAddress(address string) Address {
	return Address{
		Address:  address,
		Supernet: []string{},
		Hypernet: []string{},
	}
}

// AddSupernet adds a supernet section to the address.
func (a *Address) AddSupernet(segment string) {
	a.Supernet = append(a.Supernet, segment)
}

// AddHypernet adds a hypernet section to the address.
func (a *Address) AddHypernet(segment string) {
	a.Hypernet = append(a.Hypernet, segment)
}

// SupportsTLS determines whether the address supports TLS.
func (a Address) SupportsTLS() bool {
	// Make sure the Hypernet contains no ABBA sequence.
	for _, segment := range a.Hypernet {
		if HasAbba(segment) {
			return false
		}
	}

	// Make sure the Supernet now does contain an ABBA sequence.
	for _, segment := range a.Supernet {
		if HasAbba(segment) {
			return true
		}
	}

	return false
}

// SupportsSSL determines whether the address supports SSL.
func (a Address) SupportsSSL() bool {
	// Get the ABA/BAB triplets of the supernet and hypernet.
	supernet := FindAbaBab(a.Supernet)
	hypernet := FindAbaBab(a.Hypernet)

	// Compare and contrast in O(N^2) time.
	for s, _ := range supernet {
		for h, _ := range hypernet {
			if IsAbaBab(s, h) {
				return true
			}
		}
	}
	return false
}

// HasAbba determines whether a sequence of characters has an ABBA.
func HasAbba(input string) bool {
	// Strings less than 4 characters can't contain an ABBA.
	if len(input) < 4 {
		return false
	}

	// Run a sliding window across it checking four characters at a time.
	abba := SlidingWindow(input, 4, func(slice string) bool {
		return slice[0] == slice[3] && slice[1] == slice[2] && slice[0] != slice[1]
	})

	return len(abba) > 0
}

// IsAbaBab compares two sequences from FindAbaBab and determines whether
// they're opposites of each other (ABA -> BAB)
func IsAbaBab(a, b string) bool {
	if len(a) != 3 || len(b) != 3 {
		panic("IsAbaBab needs a 3-character sequence to work with")
	}

	return a[0] == b[1] && a[1] == b[0] && a[1] == b[2] && a[2] == b[1]
}

// FindAbaBab finds sequences of valid ABA and BAB in a given address segment.
func FindAbaBab(inputs []string) map[string]bool {
	result := map[string]bool{}

	for _, input := range inputs {
		// Strings less than 3 characters can't contain an ABA/BAB sequence.
		if len(input) < 3 {
			return result
		}

		// Run a sliding window across it checking three characters at a time for
		// both ABA and BAB sequences. This window doesn't abort (return true),
		// because we want ALL possible sequences, not just the first one.
		SlidingWindow(input, 3, func(slice string) bool {
			if slice[0] == slice[2] && slice[0] != slice[1] {
				result[slice] = true
			}
			return false
		})
	}

	return result
}

// ParseAddresses parses address lines into Address objects.
func ParseAddresses(input []string) []Address {
	result := []Address{}

	for _, line := range input {
		addr := NewAddress(line)

		// Look for bracketed sets.
		for strings.Index(line, "[") > -1 {
			match := BracketRegexp.FindStringSubmatch(line)
			if len(match) == 0 {
				fmt.Printf("Regexp fail! %s has a bracket but the regexp didn't match it; skipping\n", line)
				break
			}

			segment := match[1]
			addr.AddHypernet(segment)

			// Remove this bracketed segment once done with it.
			line = strings.Replace(line, fmt.Sprintf(`[%s]`, segment), "|", 1)
		}

		// Retrieve the supernet segments.
		for _, supernet := range strings.Split(line, "|") {
			addr.AddSupernet(supernet)
		}

		result = append(result, addr)
	}

	return result
}

// ReadFile reads lines from the input file.
func ReadFile(filename string) []string {
	fh, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var lines []string

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}

// Debug prints a debug line when $DEBUG is true.
func Debug(tmpl string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(tmpl, a...)
	}
}
