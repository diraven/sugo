package triggers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var rootCommand = &sugo.Command{
	Trigger:     "trigger",
	RootOnly:    true,
	Description: "Allows to manipulate bot trigger for the guild.",
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Get guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			return err
		}

		// Get current trigger.
		trigger := triggers.get(sg, guild.ID)

		if trigger == "" {
			trigger = sg.Self.Mention()
		}

		// Notify user about currently set bot trigger.
		if _, err := sg.RespondSuccess(m, "", "current bot trigger is set to: "+trigger); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "set",
			Description: "Sets bot trigger to the given value.",
			Usage:       "!",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				// Get guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				// Make sure prefix does not start with < as it might cause problems with mention-based triggers.
				if strings.HasPrefix(q, "<") {
					if _, err := sg.RespondDanger(m, "", "`<` is not allowed to be part of bot trigger"); err != nil {
						return err
					}
				}

				// Set our trigger.
				if err := triggers.set(sg, guild.ID, q); err != nil {
					return err
				}

				// Notify user about success.
				if _, err := sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Trigger:     "default",
			Description: "Sets trigger to default value (bot mention).",
			//Usage:         "",
			//AllowParams: false,
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				// Get guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				// Set our trigger.
				if err := triggers.setDefault(sg, guild.ID); err != nil {
					return err
				}

				// Notify user about success.
				if _, err := sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
