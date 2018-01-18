package time

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var cmdZoneSet = &sugo.Command{
	Trigger:             "set",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	AllowParams:         true,
	Description:         "Sets user timezone.",
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Make sure query is provided.
		if q == "" {
			if _, err := sg.RespondBadCommandUsage(m, c, "", ""); err != nil {
				return err
			}
		}

		// Validate the timezone.
		if _, err := getLoc(sg, m, ""); err != nil {
			_, err = sg.RespondDanger(m, "", err.Error())
			return err
		}

		// Save value.
		if err := timezones.set(sg, m.Author.ID, q); err != nil {
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondSuccess(m, "", "your new timezone: "+q); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:             "guild",
			AllowDefaultChannel: true,
			AllowParams:         true,
			RootOnly:            true,
			Description:         "Sets guild timezone.",
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				// Make sure query is provided.
				if q == "" {
					if _, err := sg.RespondBadCommandUsage(m, c, "", ""); err != nil {
						return err
					}
				}

				// Validate the timezone.
				if _, err := getLoc(sg, m, ""); err != nil {
					_, err = sg.RespondDanger(m, "", err.Error())
					return err
				}

				// Get guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				// Save value.
				if err := timezones.set(sg, guild.ID, q); err != nil {
					return err
				}

				// Respond with the resulting time to the user.
				if _, err := sg.RespondSuccess(m, "", "guild new timezone: "+q); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
