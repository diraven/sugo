package time

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

var cmdZoneReset = &sugo.Command{
	Trigger:             "reset",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	Description:         "Resets user timezone.",
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Reset value.
		if err := timezones.reset(sg, m.Author.ID); err != nil {
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondSuccess(m, "", ""); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:             "guild",
			AllowDefaultChannel: true,
			RootOnly:            true,
			Description:         "Resets guild timezone.",
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				// Get guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				// Save value.
				if err := timezones.reset(sg, guild.ID); err != nil {
					return err
				}

				// Respond with the resulting time to the user.
				if _, err := sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
