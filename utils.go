package sugo

import (
	"strings"
	"math"
)

func FmtStringsSlice(slice []string, separator string, limit int, moreText string, noMoreText string) (result string) {
	if moreText == "" {
		moreText = "..."
	}
	// Specifies if there are more items then shown.
	var more bool
	// Calculated length.
	var length int
	// Number of the first element that did not fit in.
	var i int
	// For each item in slice:
	for _, item := range slice {
		// If we are not over the limit yet:
		if length < limit-int(math.Max(float64(len(moreText)), float64(len(noMoreText)))) {
			// Increase length.
			length = length + len(item) + len(separator)
			// Increase counter.
			i = i + 1
		} else {
			// We are over the limit.
			more = true
			break
		}
	}

	result = strings.Join(slice[0:i], separator)

	if more {
		result = result + moreText
	} else {
		result = result + noMoreText
	}

	// Join elements up to the one (and not including), which is over the limit.
	return result
}
