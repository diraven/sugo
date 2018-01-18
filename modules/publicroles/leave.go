package publicroles

import (
	"github.com/diraven/sugo"
)

var leaveCmd = &sugo.Command{
	Trigger:            "leave",
	Description:        "Removes person to the public role.",
	Usage:              "role_name_or_id",
	PermittedByDefault: true,
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		// Make sure request is not empty.
		if req.Query == "" {
			_, err = sg.RespondBadCommandUsage(req, "", "")
			return err
		}

		// Try to find user public role based on query.
		roles, err := publicRoles.findUserPublicRole(sg, req, req.Query)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, req, roles, err)
		}

		// Try to remove user role.
		err = sg.GuildMemberRoleRemove(req.Guild.ID, req.Message.Author.ID, roles[0].ID)
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// Respond about operation being successful.
		_, err = sg.RespondSuccess(req, "", "you don't have `"+roles[0].Name+"` role any more")
		return err
	},
}
