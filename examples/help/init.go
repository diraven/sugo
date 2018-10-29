package help

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

// Init initializes module on the given bot.
func Init(sg *sugo.Instance) {
	sg.AddCommand(cmd)
}

func generateHelpEmbed(req *sugo.Request, c *sugo.Command) (*discordgo.MessageEmbed, error) {
	embed := &discordgo.MessageEmbed{
		Title:       c.GetPath(),
		Description: c.Description,
		Color:       sugo.ColorInfo,
	}
	// Get subcommands triggers respecting user permissions.
	subcommandsTriggers := c.GetSubcommandsTriggers(req.Sugo, req)

	if len(c.SubCommands) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Subcommands:",
				Value: strings.Join(subcommandsTriggers, ", "),
			}, &discordgo.MessageEmbedField{
				Name:  "To get help on 'subcommand' type:",
				Value: fmt.Sprintf("help %s subcommand", c.GetPath()),
			})
	}
	return embed, nil

}

var cmd = &sugo.Command{
	Trigger:     "help",
	Description: "Shows help section for the appropriate command.",
	HasParams:   true,
	Execute: func(req *sugo.Request) error {
		var err error

		// Remove help command from the string
		req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Command.Trigger))

		if req.Query == "" {
			// No arguments, just the help itself.
			_, err = req.Sugo.Session.ChannelMessageSendEmbed(req.Channel.ID, &discordgo.MessageEmbed{
				Title:       "Available commands",
				Description: strings.Join(req.Sugo.RootCommand.GetSubcommandsTriggers(req.Sugo, req), ", "),
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
		command, err := req.Sugo.FindCommand(req, req.Query)
		if err != nil {
			return err
		}
		if command != nil {
			var embed *discordgo.MessageEmbed

			embed, err = generateHelpEmbed(req, command)
			if err != nil {
				return err
			}

			_, err = req.Sugo.Session.ChannelMessageSendEmbed(req.Channel.ID, embed)
			return err
		}
		return sugo.NewCommandNotFoundError(req)
	},
}
