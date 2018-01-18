package help

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

func generateHelpEmbed(sg *sugo.Instance, req *sugo.Request, c *sugo.Command) (*discordgo.MessageEmbed, error) {
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
	subcommandsTriggers, _ := c.GetSubcommandsTriggers(sg, req)

	if len(c.SubCommands) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Subcommands:",
				Value: strings.Join(subcommandsTriggers, ", "),
			}, &discordgo.MessageEmbedField{
				Name:  "To get help on 'subcommand' type:",
				Value: fmt.Sprintf("help %s subcommand", c.Path()),
			})
	}
	return embed, nil

}

// Module shows help section for appropriate command.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "help",
		PermittedByDefault: true,
		Description:        "Shows help section for the appropriate command.",
		Usage:              "some_command",
		AllowParams:        true,
		Execute: func(sg *sugo.Instance, req *sugo.Request) error {
			var err error

			// Remove help command from the string
			req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Command.Trigger))

			if req.Query == "" {
				// No arguments, just the help itself.
				_, err = sg.ChannelMessageSendEmbed(req.Channel.ID, &discordgo.MessageEmbed{
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
			command, err := sg.FindCommand(req, req.Query)
			if err != nil {
				return err
			}
			if command != nil {
				var embed *discordgo.MessageEmbed

				embed, err = generateHelpEmbed(sg, req, command)
				if err != nil {
					return err
				}

				_, err = sg.ChannelMessageSendEmbed(req.Channel.ID, embed)
				return err
			}
			_, err = sg.RespondWarning(req, "", "I know nothing about this command, sorry...")
			return err
		},
	},
}
