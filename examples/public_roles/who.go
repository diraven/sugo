package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"sort"
)

var whoCmd = &sugo.Command{
	Trigger:            "who",
	Description:        "Lists people that have public role specified.",
	PermittedByDefault: true,
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure our role list is in sync with the server.
		storage.syncPublicRoles(sg, m)

		// Get a guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return
		}

		// Try to find role based on query.
		roles, err := storage.findGuildPublicRole(sg, m, q)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, m, roles, err)
		}

		// There is only one item in the resulting map, so we get it.
		foundRoleID, foundRoleName := pickRoleFromMap(roles)

		// Make members array we will be working with.
		memberNames := []string{}
		for _, member := range guild.Members {
			for _, roleID := range member.Roles {
				if roleID == foundRoleID {
					memberNames = append(memberNames, member.User.Username+"#"+member.User.Discriminator)
				}
			}
		}

		// Sort people.
		sort.Strings(memberNames)

		// Start building response.
		response := "people with `" + foundRoleName + "` role ```"
		response = response + sugo.FmtStringsSlice(memberNames, ", ", 1500, "...", ".")
		response = response + "```"

		_, err = sg.RespondTextMention(m, response)
		return
	},
}
