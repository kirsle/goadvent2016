package advent

import (
	"bufio"
	"os"
	"strings"
)

// ReadFile returns the lines of text in a given file, skipping blank lines.
func ReadFile(filename string) ([]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
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

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
