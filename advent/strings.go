package advent

// Truncate keeps crazy long strings from getting out of control.
func Truncate(input string, length int) string {
	if len(input) < length {
		return input
	}

	return input[0:length] + "..."
}
