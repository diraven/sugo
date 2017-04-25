package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

type BasicEmbed struct {
	Basic
	Response *discordgo.MessageEmbed
}

func (c *BasicEmbed) Execute(sg *sugo.Instance, m *discordgo.Message) (err error) {
	_, err = sg.ChannelMessageSendEmbed(m.ChannelID, c.Response)
	return
}
