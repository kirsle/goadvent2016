package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Type Room represents a parsed room.
type Room struct {
	EncryptedName string
	Sector        int
	Checksum      string
}

// Regular expressions.
var RE_RoomName = regexp.MustCompile(`^([a-z\-]+?)\-(\d+?)\[([a-z]+?)\]$`)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		os.Exit(1)
	}

	// Read the input file.
	inputLines := ReadFile(os.Args[1])

	// Parse each room.
	sectors := 0
	for _, line := range inputLines {
		room, err := ParseRoom(line)
		if err != nil {
			log.Printf("Failed to parse room %s: %v", line, err)
			continue
		}

		// Is valid?
		if err = room.Validate(); err != nil {
			Debug("Invalid room '%s': %s", room.EncryptedName, err)
			continue
		}

		// Sum up the sector ID's of the real rooms.
		sectors += room.Sector

		// Print its name.
		fmt.Printf("DECODED ROOM NAME: %s  %s  (Sector %d)\n",
			room.EncryptedName,
			room.Decrypt(),
			room.Sector,
		)
	}

	fmt.Printf("Sum of the sectors of real rooms: %d\n", sectors)
}

// ParseRoom turns a room name into a Room object.
func ParseRoom(name string) (Room, error) {
	// The parsed room object.
	room := Room{}

	// Apply the regexp first.
	match := RE_RoomName.FindStringSubmatch(name)
	if match == nil {
		return room, errors.New("Room does not match the regular expression.")
	}

	// Sector number.
	sector, err := strconv.Atoi(match[2])
	if err != nil {
		return room, errors.New("Invalid sector number")
	}

	room.EncryptedName = match[1]
	room.Sector = sector
	room.Checksum = match[3]
	return room, nil
}

// Decrypt decrypts a room name.
func (r Room) Decrypt() string {
	shift := r.Sector % 26

	// The decoded name.
	decoded := []rune{}

	// NOTE: valid letter runes range from 97 to 122 (a to z)

	for _, letter := range r.EncryptedName {
		// Dashes become spaces: easy.
		if letter == '-' {
			decoded = append(decoded, ' ')
			continue
		}

		// Shift the letter along.
		s := shiftRune(letter, shift)
		decoded = append(decoded, s)
	}

	return string(decoded)
}

// shiftRune applies a Caesar shift to a rune by an offset.
func shiftRune(r rune, shift int) rune {
	s := int(r) + shift
	if s > 'z' {
		return rune(s - 26)
	} else if s < 'a' {
		return rune(s + 26)
	}
	return rune(s)
}

// Valid validates the checksum against the encrypted name.
func (r Room) Validate() error {
	// Count the letters in the room name.
	letters := map[rune]int{}
	for _, letter := range r.EncryptedName {
		if letter == '-' {
			continue
		}

		_, ok := letters[letter]
		if !ok {
			letters[letter] = 0
		}
		letters[letter]++
	}

	// Sort the letters.
	pl := make(PairList, len(letters))
	i := 0
	for k, v := range letters {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))

	// Compute the expected checksum.
	checksum := []rune{}
	for _, letter := range pl {
		checksum = append(checksum, letter.Key)
	}

	// Checksums only have up to 5 letters.
	if len(checksum) > 5 {
		checksum = checksum[:5]
	}

	// Debugging
	Debug("Name: %v Sorted: %v [%s != %s]\n", r.EncryptedName, pl, r.Checksum, string(checksum))

	// Validate whether it matches.
	if string(checksum) != r.Checksum {
		return errors.New(fmt.Sprintf("Checksum mismatch: expected %s, got %s",
			string(checksum),
			r.Checksum,
		))
	}
	return nil
}

// Type Pair represents a character and count pair for sorting room names.
type Pair struct {
	Key   rune
	Value int
}

// Type PairList is a list of pairs to sort.
type PairList []Pair

func (slice PairList) Len() int {
	return len(slice)
}

func (slice PairList) Less(i, j int) bool {
	if slice[i].Value == slice[j].Value {
		return slice[i].Key > slice[j].Key
	}
	return slice[i].Value < slice[j].Value
}

func (slice PairList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ReadFile slurps the lines of text from a file.
func ReadFile(file string) []string {
	fh, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	lines := []string{}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

// Debug prints a debug message when $DEBUG=1.
func Debug(tmpl string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(tmpl, a...)
	}
}
