package time

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"time"
)

var cmdNow = &sugo.Command{
	Trigger:             "now",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	AllowParams:         true,
	Description:         "Checks what time is it in given timezone.",
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Get location.
		loc, err := getLoc(sg, m, q)
		if err != nil {
			if _, err := sg.RespondDanger(m, "", err.Error()); err != nil {
				return err
			}
			return nil
		}

		// Build resulting time.
		t := time.Now().In(loc)

		// Respond with the resulting time to the user.
		if _, err := sg.RespondInfo(m, "", t.Format(defaultFormat)); err != nil {
			return err
		}

		return nil
	},
}
