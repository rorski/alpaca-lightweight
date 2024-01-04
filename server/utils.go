package server

import (
	"slices"
	"strconv"
	"strings"
)

// FormatBigNum is a template filter for parsing a large number (e.g., for total equity volume)
// and making it more human readable by adding commas to it
func FormatBigNum(num uint64) string {
	var result []string
	formatNum := strconv.FormatUint(num, 10)
	if num < 1000 {
		return formatNum
	}

	splitNum := strings.Split(formatNum, "")

	// iterate backwards through the split string, since we want to add commas from right to left
	counter := 1
	for i := len(splitNum) - 1; i >= 0; i-- {
		result = append(result, string(splitNum[i]))
		if counter == 3 && i != 0 {
			// add a comma if we've reached 3 digits, and reset the counter
			result = append(result, ",")
			counter = 1
		} else {
			counter += 1
		}
	}

	// un-reverse the string and return it
	slices.Reverse(result)
	return strings.Join(result, "")
}
