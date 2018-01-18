package time

import (
	"github.com/diraven/sugo"
	"time"
)

var cmdNow = &sugo.Command{
	Trigger:             "now",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	AllowParams:         true,
	Description:         "Checks what time is it in given timezone.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Get location.
		loc, err := getLoc(sg, req)
		if err != nil {
			if _, err := sg.RespondDanger(req, "", err.Error()); err != nil {
				return err
			}
			return nil
		}

		// Build resulting time.
		t := time.Now().In(loc)

		// Respond with the resulting time to the user.
		if _, err := sg.RespondInfo(req, "", t.Format(defaultFormat)); err != nil {
			return err
		}

		return nil
	},
}
