package advent

import (
	"fmt"
	"os"
)

// Debug prints a debug message to the console when the environment variable
// `$DEBUG` is set.
func Debug(tmpl string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(tmpl, a...)
	}
}
