package help

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"strings"
	"github.com/diraven/sugo"
)

// Help shows help section for appropriate command.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "help",
		PermittedByDefault: true,
		Description:        "Shows help section for the appropriate command.",
		Usage:              "some_command",
		ParamsAllowed:      true,
		Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
			var err error

			// Remove help command from the string
			q = strings.TrimSpace(strings.TrimPrefix(q, c.Trigger))

			if q == "" {
				// No arguments, just the help itself.
				_, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       "Available commands",
					Description: strings.Join(sg.GetTriggers(), ", "),
					Color:       sugo.ColorInfo,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "To get more info on 'something' try:",
							Value: "help something",
						},
					},
				})
				return err
			}
			// Search for applicable command.
			command, err := sg.FindCommand(m, q)
			if err != nil {
				return err
			}
			if command != nil {
				var embed *discordgo.MessageEmbed
				embed, err = sg.HelpEmbed(command, m)
				if err != nil {
					return err
				}
				_, err = sg.ChannelMessageSendEmbed(m.ChannelID, embed)
				return err
			}
			_, err = sg.RespondWarning(m, "", "I know nothing about this command, sorry...")
			return err
		},
	},
}
