package time

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"time"
)

var rootCommand = &sugo.Command{
	Trigger:             "time",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	ParamsAllowed:       true,
	Description:         "Time-related tools.",
	//Usage:               defaultFormat,
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		_, err = sg.RespondBadCommandUsage(m, c, "", "")
		return err

		// TODO Finish time conversion implementation.

		// Check query.
		if q == "" {
			_, err := sg.RespondBadCommandUsage(m, c, "", "")
			return err
		}

		// Get location we are interested in.
		var loc *time.Location
		loc, err = getLoc(sg, m, "")
		if err != nil {
			if _, err := sg.RespondDanger(m, "", err.Error()); err != nil {
				return err
			}
		}

		// Try to parse time.
		var t time.Time
		if t, _ = time.ParseInLocation(defaultFormat, q, loc); err != nil {
			_, err := sg.RespondDanger(m, "", "unable to parse time, the format should be: "+defaultFormat)
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondInfo(m, "", t.In(loc).Format(defaultFormat)); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		cmdNow,
		cmdZone,
	},
}
