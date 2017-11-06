package altaxi

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

// Command contains all ed-related stuff.
var orderCancel = &sugo.Command{
	Trigger:            "cancel",
	PermittedByDefault: true,
	Description:        "Cancels last order.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		_, exists := pendingOrders[m.Author.ID]
		if exists {
			delete(pendingOrders, m.Author.ID)
			_, err = sg.RespondSuccessMention(m, "")
			if err != nil {
				return
			}
		} else {
			_, err = sg.RespondTextMention(m, "You have no pending orders.")
		}
		return
	},
}
