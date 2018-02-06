package publicroles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var delCmd = &sugo.Command{
	Trigger:             "del",
	Description:         "Makes given role not public (does not delete the role itself).",
	PermissionsRequired: discordgo.PermissionManageRoles,
	HasParams:           true,
	Execute: func(req *sugo.Request) error {
		var err error

		// Make sure query is not empty.
		if req.Query == "" {
			_, err = req.RespondBadCommandUsage("", "")
			return err
		}

		// Try to find role based on query.
		roles, err := storage.findGuildPublicRole(req, req.Query)
		if err != nil {
			return respondFuzzyRolesSearchIssue(req, roles, err)
		}

		// Delete role.
		err = storage.del(roles[0].ID)
		if err != nil {
			return err
		}

		// Notify user about success of the operation.
		_, err = req.RespondSuccess("", "role `"+roles[0].Name+"` is not public any more")
		return err
	},
}
