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
	ParamsAllowed: true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Try to find role based on query.
		roles, err := publicRoles.findUserPublicRole(sg, m, q)

		// Start building response.
		var response string

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response = response + "```\n"
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
