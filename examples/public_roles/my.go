package public_roles

import (
	"context"
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

var myCmd = &sugo.Command{
	Trigger:            "my",
	Description:        "Lists public roles you are in.",
	PermittedByDefault: true,
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure our role list is in sync with the server.
		storage.syncPublicRoles(sg, m)

		// Try to find role based on query.
		roles, err := storage.findUserPublicRole(sg, m, q)

		// Start building response.
		var response string

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response = response + "```"
			response = response + sugo.FmtStringsSlice(rolesToRoleNames(roles), "\n", 1500, "\n...", "")
			response = response + "```"
			_, err = sg.RespondTextMention(m, response)
		} else {
			if q == "" {
				_, err = sg.RespondTextMention(m, "you got no roles")
			} else {
				_, err = sg.RespondTextMention(m, "you got no such roles")
			}
		}

		return err
	},
}
