package aliases

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
	"strings"
)

var rootCommand = &sugo.Command{
	Trigger:     "aliases",
	RootOnly:    true,
	Description: "Allows to manipulate aliases. Lists all aliases.",
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

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
			Trigger:       "set",
			Description:   "Adds new or updates existent alias.",
			Usage:         "some_alias -> command [subcommand ...]",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				ss := strings.Split(q, "->")
				if len(ss) < 2 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return err
				}
				alias := strings.TrimSpace(ss[0])
				commandPath := strings.TrimSpace(ss[1])

				// Try to find command.
				command, err := sg.FindCommand(m, commandPath)
				if command == nil {
					_, err = sg.RespondCommandNotFound(m)
					return err
				}

				guild := ctx.Value(sugo.CtxKey("guild")).(*discordgo.Guild)

				aliases.set(sg, guild, alias, commandPath)
				_, err = sg.RespondSuccess(m, "")
				return err
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified alias.",
			Usage:       "some_alias",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				guild := ctx.Value(sugo.CtxKey("guild")).(*discordgo.Guild)
				alias := aliases.get(sg, guild, q)

				if alias == "" {
					_, err = sg.RespondDanger(m, "Alias \""+q+"\" was not found.")
					return err
				}

				aliases.del(sg, guild, alias)
				_, err = sg.RespondSuccess(m, "")
				return err
			},
		},
		{
			Trigger:     "swap",
			Description: "Swaps specified shortcuts.",
			Usage:       "1 2",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				guild := ctx.Value(sugo.CtxKey("guild")).(*discordgo.Guild)

				ss := strings.Split(q, " ")
				if len(ss) < 2 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return err
				}

				alias1 := aliases.get(sg, guild, ss[0])
				if alias1 == "" {
					sg.RespondDanger(m, "alias \""+ss[0]+"\" not found")
					return err
				}

				alias2 := aliases.get(sg, guild, ss[1])
				if alias2 == "" {
					sg.RespondDanger(m, "alias \""+ss[1]+"\" not found")
					return err
				}

				aliases.swap(sg, guild, alias1, alias2)

				_, err = sg.RespondSuccess(m, "")
				return err
			},
		},
	},
}