package helpers

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
)

func UserAsMention(u discordgo.User) (s string, err error) {
	s = fmt.Sprintf("<@%s>", u.ID)
	return
}

func ChannelAsMention(c discordgo.Channel) (s string, err error) {
	s = fmt.Sprintf("<#%s>", c.ID)
	return
}

func TimeToDiscordFormat(t time.Time) (s string) {
	s = t.Format(time.RFC3339)
	return
}