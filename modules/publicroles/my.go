package publicroles

import (
	"github.com/diraven/sugo"
)

var myCmd = &sugo.Command{
	Trigger:            "my",
	Description:        "Lists public roles you are in.",
	PermittedByDefault: true,
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Try to find role based on query.
		roles, err := publicRoles.findUserPublicRole(sg, req, req.Query)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := "public roles you are in:\n\n"
			response = response + sugo.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
			_, err = sg.RespondInfo(req, "", response)
		} else {
			if req.Query == "" {
				_, err = sg.RespondWarning(req, "", "you got no roles")
			} else {
				_, err = sg.RespondWarning(req, "", "you got no such roles")
			}
		}

		return err
	},
}
