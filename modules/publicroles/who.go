package publicroles

import (
	"github.com/diraven/sugo"
)

var whoCmd = &sugo.Command{
	Trigger:            "who",
	Description:        "Lists people that have public role specified.",
	PermittedByDefault: true,
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, req, req.Query)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, req, roles, err)
		}

		// Make members array we will be working with.
		var memberMentions []string
		for _, member := range req.Guild.Members {
			for _, roleID := range member.Roles {
				if roleID == roles[0].ID {
					memberMentions = append(memberMentions, member.User.Mention())
				}
			}
		}

		// Start building response.
		response := sugo.FmtStringsSlice(memberMentions, ", ", "", 1500, "...", "")

		_, err = sg.RespondInfo(req, "", "people with `"+roles[0].Name+"` role:\n\n"+response)
		return err
	},
}
