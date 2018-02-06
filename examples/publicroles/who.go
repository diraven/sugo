package publicroles

import (
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
)

var whoCmd = &sugo.Command{
	Trigger:     "who",
	Description: "Lists people that have public role specified.",
	HasParams:   true,
	Execute: func(req *sugo.Request) error {
		var err error

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

		// Make members array we will be working with.
		var memberMentions []string
		for _, member := range guild.Members {
			for _, roleID := range member.Roles {
				if roleID == roles[0].ID {
					memberMentions = append(memberMentions, member.User.Mention())
				}
			}
		}

		// Start building response.
		response := helpers.FmtStringsSlice(memberMentions, ", ", "", 1500, "...", "")

		_, err = req.RespondInfo("", "people with `"+roles[0].Name+"` role:\n\n"+response)
		return err
	},
}
