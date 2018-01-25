package helpers

import (
	"time"
	"math"
	"strings"
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

	for i := range slice {
		slice[i] = decorator + slice[i] + decorator
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

//func ConsumePrefix(s string, p string) string {
//	if strings.HasPrefix(s, p) {
//		s = strings.TrimPrefix(s, p)
//		return strings.TrimSpace(s)
//	}
//	return s
//}
//
//// PointerToString returns pointer to the string given.
//func PointerToString(inputStr string) (pointerToStr *string) {
//	return &inputStr
//}
//
//// GetRoleById returns discordgo.Role based on role ID provided or nil if role not found
//func GetRoleById(guild *discordgo.Guild, roleID string) *discordgo.Role {
//	for _, role := range guild.Roles {
//		if role.ID == roleID {
//			return role
//		}
//	}
//	return nil
//}
//
//// DiscordIDCreationTime gets the time of creation of any ID.
//func DiscordIDCreationTime(ID string) (time.Time, error) {
//	i, err := strconv.ParseInt(ID, 10, 64)
//	if err != nil {
//		return
//	}
//	timestamp := (i >> 22) + 1420070400000
//	t = time.Unix(timestamp/1000, 0)
//	return
//}
//
//// BoolToInt converts bool:true to int:1 and bool:false to int:0.
//func BoolToInt(input bool) (output int) {
//	if input {
//		return 1
//	}
//	return 0
//}
