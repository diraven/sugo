package publicroles

import (
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
)

var myCmd = &sugo.Command{
	Trigger:     "my",
	Description: "Lists public roles you are in.",
	Execute: func(req *sugo.Request) error {
		// Try to find role based on query.
		roles, err := storage.findUserPublicRole(req, req.Query)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := "public roles you are in:\n\n"
			response = response + helpers.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
			_, err = req.RespondInfo("", response)
		} else {
			if req.Query == "" {
				_, err = req.RespondWarning("", "you got no roles")
			} else {
				_, err = req.RespondWarning("", "you got no such roles")
			}
		}

		return err
	},
}
