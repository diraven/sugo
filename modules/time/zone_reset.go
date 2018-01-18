package time

import (
	"github.com/diraven/sugo"
)

var cmdZoneReset = &sugo.Command{
	Trigger:             "reset",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	Description:         "Resets user timezone.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Reset value.
		if err := timezones.reset(sg, req.Message.Author.ID); err != nil {
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondSuccess(req, "", ""); err != nil {
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
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				// Save value.
				if err := timezones.reset(sg, req.Guild.ID); err != nil {
					return err
				}

				// Respond with the resulting time to the user.
				if _, err := sg.RespondSuccess(req, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
