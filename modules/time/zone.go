package time

import (
	"github.com/diraven/sugo"
)

var cmdZone = &sugo.Command{
	Trigger:             "zone",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	Description:         "Shows user timezone.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Get user timezone value.
		tzone, err := timezones.get(sg, req)
		if err != nil {
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondInfo(req, "", "your timezone: "+tzone); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		cmdZoneSet,
		cmdZoneReset,
	},
}
