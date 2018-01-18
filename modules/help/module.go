package help

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

func generateHelpEmbed(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message) (*discordgo.MessageEmbed, error) {
	embed := &discordgo.MessageEmbed{
		Title:       c.Path(),
		Description: c.Description,
		Color:       sugo.ColorInfo,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Usage:",
				Value: c.Path() + " " + c.Usage,
			},
		},
	}
	// Get subcommands triggers respecting user permissions.
	subcommandsTriggers, _ := c.GetSubcommandsTriggers(sg, m)

	if len(c.SubCommands) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Subcommands:",
				Value: strings.Join(subcommandsTriggers, ", "),
			}, &discordgo.MessageEmbedField{
				Name:  "To get help on 'subcommand' type:",
				Value: fmt.Sprintf("help %s subcommand", c.Trigger),
			})
	}
	return embed, nil

}

// Help shows help section for appropriate command.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "help",
		PermittedByDefault: true,
		Description:        "Shows help section for the appropriate command.",
		Usage:              "some_command",
		AllowParams:        true,
		Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
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

				embed, err = generateHelpEmbed(sg, command, m)
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
