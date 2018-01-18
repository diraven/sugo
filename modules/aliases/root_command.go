package aliases

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var rootCommand = &sugo.Command{
	Trigger:     "aliases",
	RootOnly:    true,
	Description: "Allows to manipulate aliases. Lists all aliases.",
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			return err
		}

		var result string
		for alias, commandPath := range *aliases.all(guild) {
			result = result + alias + " -> " + commandPath + "\n"
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Currently configured aliases",
			Description: result,
		}

		sg.ChannelMessageSendEmbed(m.ChannelID, embed)
		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "set",
			Description: "Adds new or updates existent alias.",
			Usage:       "some_alias -> command [subcommand ...]",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				ss := strings.Split(q, "->")
				if len(ss) < 2 {
					_, err := sg.RespondBadCommandUsage(m, c, "", "")
					return err
				}
				alias := strings.TrimSpace(ss[0])
				commandPath := strings.TrimSpace(ss[1])

				// Try to find command.
				command, err := sg.FindCommand(m, commandPath)
				if err != nil {
					return err
				}
				if command == nil {
					_, err := sg.RespondCommandNotFound(m)
					return err
				}

				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				aliases.set(sg, guild, alias, commandPath)
				if _, err := sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified alias.",
			Usage:       "some_alias",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				alias := aliases.get(sg, guild, q)

				if alias == "" {
					_, err := sg.RespondDanger(m, "", "Alias \""+q+"\" was not found.")
					return err
				}

				aliases.del(sg, guild, alias)
				if _, err := sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Trigger:     "swap",
			Description: "Swaps specified shortcuts.",
			Usage:       "1 2",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				ss := strings.Split(q, " ")
				if len(ss) < 2 {
					_, err := sg.RespondBadCommandUsage(m, c, "", "")
					return err
				}

				alias1 := aliases.get(sg, guild, ss[0])
				if alias1 == "" {
					_, err := sg.RespondDanger(m, "", "alias `"+ss[0]+"` not found")
					return err
				}

				alias2 := aliases.get(sg, guild, ss[1])
				if alias2 == "" {
					_, err := sg.RespondDanger(m, "", "alias `"+ss[1]+"` not found")
					return err
				}

				aliases.swap(sg, guild, alias1, alias2)

				if _, err := sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
