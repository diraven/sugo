package helpers

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
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

func PointerToString(input_str string)(pointer_to_str *string) {
	return &input_str
}