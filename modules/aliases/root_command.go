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
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var result string
		for alias, commandPath := range *aliases.all(req.Guild) {
			result = result + alias + " -> " + commandPath + "\n"
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Currently configured aliases",
			Description: result,
		}

		sg.ChannelMessageSendEmbed(req.Message.ChannelID, embed)
		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "set",
			Description: "Adds new or updates existent alias.",
			Usage:       "some_alias -> command [subcommand ...]",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				ss := strings.Split(req.Query, "->")
				if len(ss) < 2 {
					_, err := sg.RespondBadCommandUsage(req, "", "")
					return err
				}
				alias := strings.TrimSpace(ss[0])
				commandPath := strings.TrimSpace(ss[1])

				// Try to find command.
				command, err := sg.FindCommand(req, commandPath)
				if err != nil {
					return err
				}
				if command == nil {
					_, err := sg.RespondCommandNotFound(req)
					return err
				}

				aliases.set(sg, req.Guild, alias, commandPath)
				if _, err := sg.RespondSuccess(req, "", ""); err != nil {
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
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				alias := aliases.get(sg, req.Guild, req.Query)

				if alias == "" {
					_, err := sg.RespondDanger(req, "", "Alias \""+req.Query+"\" was not found.")
					return err
				}

				aliases.del(sg, req.Guild, alias)
				if _, err := sg.RespondSuccess(req, "", ""); err != nil {
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
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				ss := strings.Split(req.Query, " ")
				if len(ss) < 2 {
					_, err := sg.RespondBadCommandUsage(req, "", "")
					return err
				}

				alias1 := aliases.get(sg, req.Guild, ss[0])
				if alias1 == "" {
					_, err := sg.RespondDanger(req, "", "alias `"+ss[0]+"` not found")
					return err
				}

				alias2 := aliases.get(sg, req.Guild, ss[1])
				if alias2 == "" {
					_, err := sg.RespondDanger(req, "", "alias `"+ss[1]+"` not found")
					return err
				}

				aliases.swap(sg, req.Guild, alias1, alias2)

				if _, err := sg.RespondSuccess(req, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
