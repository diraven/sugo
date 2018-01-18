package publicroles

import (
	"github.com/diraven/sugo"
)

var joinCmd = &sugo.Command{
	Trigger:            "join",
	Description:        "Joins person to the public role.",
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

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, req, req.Query)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, req, roles, err)
		}

		// Try to assign role.
		err = sg.GuildMemberRoleAdd(req.Guild.ID, req.Message.Author.ID, roles[0].ID)
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// Respond about successful operation.
		_, err = sg.RespondSuccess(req, "", "you now have `"+roles[0].Name+"` role")
		return err
	},
}
