package publicroles

import (
	"github.com/diraven/sugo"
)

var rootCommand = &sugo.Command{
	Trigger:            "public_roles",
	Description:        "Allows to manipulate public roles.",
	PermittedByDefault: true,
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, req, req.Query)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := sugo.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
			_, err = sg.RespondInfo(req, "available public roles", response)
		} else {
			_, err = sg.RespondDanger(req, "", "nothing found")
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
