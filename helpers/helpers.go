package helpers

import (
	"time"
	"math"
	"strings"
	"strconv"
)

// TimeToDiscordTimestamp returns time in a format that Discord API accepts.
func TimeToDiscordTimestamp(t time.Time) (s string) {
	s = t.Format(time.RFC3339)
	return
}

// DiscordTimestampToTime returns time.Time parsed from discord API time string.
func DiscordTimestampToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// FmtStringsSlice formats strings slice according to the given parameters.
func FmtStringsSlice(slice []string, separator string, decorator string, limit int, moreText string, noMoreText string) (result string) {
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
			length = length + len(item) + len(separator) + 2*len(decorator)
			// Increase counter.
			i = i + 1
		} else {
			// We are over the limit.
			more = true
			break
		}
	}

	// Decorate slice with decorators.
	for i := range slice {
		slice[i] = decorator + slice[i] + decorator
	}

	// Join elements up to the one (and not including), which is over the limit.
	result = strings.Join(slice[0:i], separator)

	// Add moreText if we have cut some elements off.
	if more {
		result = result + moreText
	} else {
		result = result + noMoreText
	}

	return result
}

// DiscordIDCreationTime gets the time of creation of any ID.
func DiscordIDCreationTime(ID string) (*time.Time, error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return nil, err
	}
	timestamp := (i >> 22) + 1420070400000
	t := time.Unix(timestamp/1000, 0)
	return &t, nil
}
