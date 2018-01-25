package publicroles

import (
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
)

var myCmd = &sugo.Command{
	Trigger:            "my",
	Description:        "Lists public roles you are in.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Try to find role based on query.
		roles, err := storage.findUserPublicRole(sg, req, req.Query)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := "public roles you are in:\n\n"
			response = response + helpers.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
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
