package time

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
)

var cmdZone = &sugo.Command{
	Trigger:             "zone",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	Description:         "Shows user timezone.",
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Get user timezone value.
		tzone, err := timezones.get(sg, m)
		if err != nil {
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondInfo(m, "", "your timezone: "+tzone); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		cmdZoneSet,
		cmdZoneReset,
	},
}
