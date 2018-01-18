package publicroles

import (
	"github.com/diraven/sugo"
)

var delCmd = &sugo.Command{
	Trigger:     "del",
	Description: "Makes given role not public (does not delete the role itself).",
	Usage:       "role_name_or_id",
	AllowParams: true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		// Make sure query is not empty.
		if req.Query == "" {
			_, err = sg.RespondBadCommandUsage(req, "", "")
			return err
		}

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, req, req.Query)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, req, roles, err)
		}

		// Delete role.
		publicRoles.del(sg, req.Guild.ID, roles[0].ID)

		// Notify user about success of the operation.
		_, err = sg.RespondSuccess(req, "", "role `"+roles[0].Name+"` is not public any more")
		return err
	},
}
