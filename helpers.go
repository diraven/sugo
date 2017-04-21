package sugo

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

func UserAsMention(u discordgo.User) (s string, err error) {
	s = fmt.Sprintf("<@%s>", u.ID)
	return
}

func ChannelAsMention(c discordgo.Channel) (s string, err error) {
	s = fmt.Sprintf("<#%s>", c.ID)
	return
}