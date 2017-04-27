package helpers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

// UserAsMention returns a string in a form of <@0000000000> to be included into discord message.
func UserAsMention(u *discordgo.User) (s string) {
	s = fmt.Sprintf("<@%s>", u.ID)
	return
}

// ChannelAsMention returns a string in a form of <#0000000000> to be included into discord message.
func ChannelAsMention(c *discordgo.Channel) (s string) {
	s = fmt.Sprintf("<#%s>", c.ID)
	return
}

// TimeToDiscordTimestamp returns time in a format that Discord API accepts.
func TimeToDiscordTimestamp(t time.Time) (s string) {
	s = t.Format(time.RFC3339)
	return
}

// DiscordTimestampToTime returns time.Time parsed from discord API time string.
func DiscordTimestampToTime(s string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return
	}
	return
}

// PointerToString returns pointer to the string given.
func PointerToString(inputStr string) (pointerToStr *string) {
	return &inputStr
}

// DiscordIDCreationTime gets the time of creation of any ID.
func DiscordIDCreationTime(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}

// BoolToInt converts bool:true to int:1 and bool:false to int:0.
func BoolToInt(input bool) (output int) {
	if input {
		return 1
	}
	return 0
}

// ConsumeTerm removes and returns first term from the given string.
//func ConsumeTerm(input *string) (term string) {
//	*input = strings.TrimSpace(*input)
//	nextSpaceIndex := strings.Index(*input, " ")
//	if nextSpaceIndex < 0 {
//		term = *input
//		*input = ""
//	} else {
//		term = (*input)[:nextSpaceIndex]
//		*input = (*input)[nextSpaceIndex+1:]
//	}
//	return
//}
