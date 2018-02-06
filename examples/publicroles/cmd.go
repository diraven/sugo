package publicroles

import (
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
)

var cmd = &sugo.Command{
	Trigger:     "public_roles",
	Description: "Allows to manipulate public roles.",
	Execute: func(req *sugo.Request) error {
		// Try to find role based on query.
		roles, err := storage.findGuildPublicRole(req, req.Query)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := helpers.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
			_, err = req.RespondInfo("available public roles", response)
		} else {
			_, err = req.RespondDanger("", "nothing found")
		}

		return err
	},
	SubCommands: []*sugo.Command{
		myCmd,
		whoCmd,
		addCmd,
		delCmd,
		joinCmd,
		leaveCmd,
		createCmd,
		statsCmd,
	},
}
