package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Type Password contains the slots for the password.
const PasswordLength = 8

type Password struct {
	Code   [PasswordLength]rune // The actual characters of the password.
	Filled [PasswordLength]bool // Layer mask for which runes we've unlocked.
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <puzzle input>")
		os.Exit(1)
	}

	input := os.Args[1]
	index := -1
	password := Password{}

	// Search for that password.
	for !password.Cracked() {
		index++
		hash := Hash(input + strconv.Itoa(index))
		if strings.HasPrefix(hash, "00000") {
			// Interesting hash found!
			position := hash[5]
			value := hash[6]
			if position < '0' || position > '7' {
				continue
			}

			// Turn the position symbol into a normal int. Guaranteed to work,
			// since we excluded impossible positions just above.
			pos, _ := strconv.Atoi(string(position))
			Debug("Index=%d hash=%s position=%s value=%s\n", index, hash, pos, string(value))
			password.Fill(pos, rune(value))
		}
	}

	fmt.Printf("The password is: %s\n", password.String())
}

// Fill enters a password symbol, if that position wasn't already found.
func (p *Password) Fill(position int, value rune) {
	if !p.Filled[position] {
		p.Code[position] = rune(value)
		p.Filled[position] = true
		fmt.Printf("Cracking: %s\n", p.String())
	}
}

// String returns a string version of the password.
func (p *Password) String() string {
	result := []string{}
	for i, _ := range p.Code {
		if p.Filled[i] {
			result = append(result, string(p.Code[i]))
		} else {
			result = append(result, "-")
		}
	}
	return strings.Join(result, "")
}

// Cracked tests whether the password was fully cracked.
func (p *Password) Cracked() bool {
	for _, found := range p.Filled {
		if !found {
			return false
		}
	}
	return true
}

// Hash creates an MD5 hash of a string and returns hexadecimal.
func Hash(in string) string {
	hasher := md5.New()
	hasher.Write([]byte(in))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Debug prints a debug message when $DEBUG=1.
func Debug(tmpl string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(tmpl, a...)
	}
}
