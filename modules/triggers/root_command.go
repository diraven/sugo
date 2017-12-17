package triggers

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
	"strings"
)

var rootCommand = &sugo.Command{
	Trigger:     "trigger",
	RootOnly:    true,
	Description: "Allows to manipulate bot trigger for the guild.",
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Get guild.
		guild := ctx.Value(sugo.CtxKey("guild")).(*discordgo.Guild)

		// Get current trigger.
		trigger := triggers.get(sg, guild.ID)

		if trigger == "" {
			trigger = sg.Self.Mention()
		}

		// Notify user about currently set bot trigger.
		if _, err := sg.RespondSuccess(m, "current bot trigger is set to: +trigger"); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:       "set",
			Description:   "Sets bot trigger to the given value.",
			Usage:         "!",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				// Get guild.
				guild := ctx.Value(sugo.CtxKey("guild")).(*discordgo.Guild)

				// Make sure prefix does not start with < as it might cause problems with mention-based triggers.
				if strings.HasPrefix(q, "<") {
					if _, err := sg.RespondDanger(m, "`<` is not allowed to be part of bot trigger"); err != nil {
						return err
					}
				}

				// Set our trigger.
				if err := triggers.set(sg, guild.ID, q); err != nil {
					return err
				}

				// Notify user about success.
				if _, err := sg.RespondSuccess(m, ""); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Trigger:     "default",
			Description: "Sets trigger to default value (bot mention).",
			//Usage:         "",
			//ParamsAllowed: false,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				// Get guild.
				guild := ctx.Value(sugo.CtxKey("guild")).(*discordgo.Guild)

				// Set our trigger.
				if err := triggers.setDefault(sg, guild.ID); err != nil {
					return err
				}

				// Notify user about success.
				if _, err := sg.RespondSuccess(m, ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
