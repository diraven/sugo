package public_roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var myCmd = &sugo.Command{
	Trigger:            "my",
	Description:        "Lists public roles you are in.",
	PermittedByDefault: true,
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Try to find role based on query.
		roles, err := publicRoles.findUserPublicRole(sg, m, q)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := "public roles you are in:\n\n"
			response = response + sugo.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
			_, err = sg.RespondInfo(m, "", response)
		} else {
			if q == "" {
				_, err = sg.RespondWarning(m, "", "you got no roles")
			} else {
				_, err = sg.RespondWarning(m, "", "you got no such roles")
			}
		}

		return err
	},
}
