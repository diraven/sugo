package time

import (
	"github.com/diraven/sugo"
)

var cmdZoneSet = &sugo.Command{
	Trigger:             "set",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	AllowParams:         true,
	Description:         "Sets user timezone.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Make sure query is provided.
		if req.Query == "" {
			if _, err := sg.RespondBadCommandUsage(req, "", ""); err != nil {
				return err
			}
		}

		// Validate the timezone.
		if _, err := getLoc(sg, req); err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// Save value.
		if err := timezones.set(sg, req.Message.Author.ID, req.Query); err != nil {
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondSuccess(req, "", "your new timezone: "+req.Query); err != nil {
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
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				// Make sure query is provided.
				if req.Query == "" {
					if _, err := sg.RespondBadCommandUsage(req, "", ""); err != nil {
						return err
					}
				}

				// Validate the timezone.
				if _, err := getLoc(sg, req); err != nil {
					_, err = sg.RespondDanger(req, "", err.Error())
					return err
				}

				// Save value.
				if err := timezones.set(sg, req.Guild.ID, req.Query); err != nil {
					return err
				}

				// Respond with the resulting time to the user.
				if _, err := sg.RespondSuccess(req, "", "guild new timezone: "+req.Query); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
