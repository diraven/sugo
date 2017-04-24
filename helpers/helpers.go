package helpers

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
	"strconv"
	"strings"
)

func UserAsMention(u *discordgo.User) (s string) {
	s = fmt.Sprintf("<@%s>", u.ID)
	return
}

func ChannelAsMention(c *discordgo.Channel) (s string) {
	s = fmt.Sprintf("<#%s>", c.ID)
	return
}

func TimeToDiscordTimestamp(t time.Time) (s string) {
	s = t.Format(time.RFC3339)
	return
}

func DiscordTimestampToTime(s string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return
	}
	return
}

func PointerToString(input_str string) (pointer_to_str *string) {
	return &input_str
}

// GetCreationTime is a helper function to get the time of creation of any ID.
// ID: ID to get the time from
func DiscordIDCreationTime(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}

// ConsumeTerm removes and returns first term from the given string.
func ConsumeTerm(input *string) (term string) {
	*input = strings.TrimSpace(*input)
	nextSpaceIndex := strings.Index(*input, " ")
	if nextSpaceIndex < 0 {
		term = *input
		*input = ""
	} else {
		term = (*input)[:nextSpaceIndex]
		*input = (*input)[nextSpaceIndex+1:]
	}
	return
}
