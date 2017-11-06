package altaxi

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var cmdReportReset = &sugo.Command{
	Trigger:            "reset",
	PermittedByDefault: false,
	Description:        "Resets all daily earnings.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		err = postReport(ctx, c, q, sg, m)
		if err != nil {
			return
		}

		reportStorage.Items = make(map[string]reportItem)
		_, err = sg.RespondSuccessMention(m, "")

		return
	},
}
