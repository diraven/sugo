package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"github.com/diraven/sugo"
)

type Basic struct {
	Triggers []string
	Response *string
}

func (c Basic) Test(sg *sugo.Instance, m *discordgo.Message) (is_applicable bool) {
	is_applicable = false
	for _, trigger := range c.Triggers {
		if strings.Contains(m.Content, trigger) {
			is_applicable = true
			return
		}
	}
	return
}

func (c Basic) Execute(sg *sugo.Instance, m *discordgo.Message) (err error) {
	_, err = sg.ChannelMessageSend(m.ChannelID, *c.Response)
	return
}
