package time

import (
	"github.com/diraven/sugo"
	"time"
)

var rootCommand = &sugo.Command{
	Trigger:             "time",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	AllowParams:         true,
	Description:         "Time-related tools.",
	//Usage:               defaultFormat,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		_, err = sg.RespondBadCommandUsage(req, "", "")
		return err

		// TODO Finish time conversion implementation.

		// Check query.
		if req.Query == "" {
			_, err := sg.RespondBadCommandUsage(req, "", "")
			return err
		}

		// Get location we are interested in.
		var loc *time.Location
		loc, err = getLoc(sg, req)
		if err != nil {
			if _, err := sg.RespondDanger(req, "", err.Error()); err != nil {
				return err
			}
		}

		// Try to parse time.
		var t time.Time
		if t, _ = time.ParseInLocation(defaultFormat, req.Query, loc); err != nil {
			_, err := sg.RespondDanger(req, "", "unable to parse time, the format should be: "+defaultFormat)
			return err
		}

		// Respond with the resulting time to the user.
		if _, err := sg.RespondInfo(req, "", t.In(loc).Format(defaultFormat)); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		cmdNow,
		cmdZone,
	},
}
