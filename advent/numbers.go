package advent

import "strconv"

// StringsToInts turns a slice of strings into a slice of ints.
func StringsToInts(input []string) ([]int, error) {
	var result []int
	for _, str := range input {
		number, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result = append(result, number)
	}
	return result, nil
}
