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

func generateHelpEmbed(sg *sugo.Instance, req *sugo.Request, c *sugo.Command) (*discordgo.MessageEmbed, error) {
	embed := &discordgo.MessageEmbed{
		Title:       c.GetPath(),
		Description: c.Description,
		Color:       sugo.ColorInfo,
	}
	// Get subcommands triggers respecting user permissions.
	subcommandsTriggers := c.GetSubcommandsTriggers(sg, req)

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
	Execute: func(sg *sugo.Instance, r *sugo.Request) error {
		var err error

		// Remove help command from the string
		r.Query = strings.TrimSpace(strings.TrimPrefix(r.Query, r.Command.Trigger))

		if r.Query == "" {
			// No arguments, just the help itself.
			_, err = sg.Session.ChannelMessageSendEmbed(r.Channel.ID, &discordgo.MessageEmbed{
				Title:       "Available commands",
				Description: strings.Join(sg.RootCommand.GetSubcommandsTriggers(sg, r), ", "),
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
		command, err := sg.FindCommand(r, r.Query)
		if err != nil {
			return err
		}
		if command != nil {
			var embed *discordgo.MessageEmbed

			embed, err = generateHelpEmbed(sg, r, command)
			if err != nil {
				return err
			}

			_, err = sg.Session.ChannelMessageSendEmbed(r.Channel.ID, embed)
			return err
		}
		_, err = sg.RespondWarning(r, "", "I know nothing about this command, sorry...")
		return err
	},
}
