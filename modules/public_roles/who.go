package public_roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var whoCmd = &sugo.Command{
	Trigger:            "who",
	Description:        "Lists people that have public role specified.",
	PermittedByDefault: true,
	ParamsAllowed:      true,
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Get a guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			_, err = sg.RespondDanger(m, "", err.Error())
			return err
		}

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, m, q)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, m, roles, err)
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
		response := sugo.FmtStringsSlice(memberMentions, ", ", "", 1500, "...", "")

		_, err = sg.RespondInfo(m, "", "people with `"+roles[0].Name+"` role:\n\n"+response)
		return err
	},
}
