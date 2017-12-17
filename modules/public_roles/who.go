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
	ParamsAllowed: true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Get a guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			_, err = sg.RespondDanger(m, err.Error())
			return err
		}

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, m, q)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, m, roles, err)
		}

		// Make members array we will be working with.
		var memberNames []string
		for _, member := range guild.Members {
			for _, roleID := range member.Roles {
				if roleID == roles[0].ID {
					memberNames = append(memberNames, member.User.Username+"#"+member.User.Discriminator)
				}
			}
		}

		// Sort people.
		sort.Strings(memberNames)

		// Start building response.
		response := "people with `" + roles[0].Name + "` role ```\n"
		response = response + sugo.FmtStringsSlice(memberNames, ", ", 1500, "...", ".")
		response = response + "```"

		_, err = sg.Respond(m, response)
		return err
	},
}
