package publicroles

import (
	"github.com/diraven/sugo"
)

var joinCmd = &sugo.Command{
	Trigger:     "join",
	Description: "Joins person to the public role.",
	HasParams:   true,
	Execute: func(req *sugo.Request) error {
		var err error

		// Make sure request is not empty.
		if req.Query == "" {
			_, err = req.RespondBadCommandUsage("", "")
			return err
		}

		// Try to find role based on query.
		roles, err := storage.findGuildPublicRole(req, req.Query)
		if err != nil {
			return respondFuzzyRolesSearchIssue(req, roles, err)
		}

		// Get guild.
		guild, err := req.GetGuild()
		if err != nil {
			return err
		}

		// Try to assign role.
		err = req.Sugo.Session.GuildMemberRoleAdd(guild.ID, req.Message.Author.ID, roles[0].ID)
		if err != nil {
			_, err = req.RespondDanger("", err.Error())
			return err
		}

		// Respond about successful operation.
		_, err = req.RespondSuccess("", "you now have `"+roles[0].Name+"` role")
		return err
	},
}
