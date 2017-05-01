package sugo

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
	"strings"
)

// Help shows help section for appropriate command.
var CmdHelp = &Command{
	Trigger:     "help",
	Permissions: []int{PermissionNone},
	Description: "Shows help section for the appropriate command.",
	Usage:       "somecommand",
	Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
		// Remove help command from the string
		q = strings.TrimSpace(strings.TrimPrefix(q, c.Trigger))

		if q == "" {
			// No arguments, just the help itself.
			_, err = sg.RespondEmbed(m, &discordgo.MessageEmbed{
				Title:       "Available commands",
				Description: strings.Join(sg.triggers(), ", "),
				Color:       ColorInfo,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "To get more info on 'something' try:",
						Value: helpers.UserAsMention(sg.Self) + " help something",
					},
				},
			})
			return
		}
		// Search for applicable command.
		command, err := findCommand(q, m, sg.commands())
		if err != nil {
			return
		}
		if command != nil {
			var embed *discordgo.MessageEmbed
			embed, err = sg.helpEmbed(command)
			if err != nil {
				return
			}
			_, err = sg.RespondEmbed(m, embed)
			return
		}
		_, err = sg.RespondTextMention(m, "I know nothing about this command, sorry...")
		return
	},
}
