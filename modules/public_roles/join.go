package public_roles

import (
	"context"
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

var joinCmd = &sugo.Command{
	Trigger:            "join",
	Description:        "Joins person to the public role.",
	Usage:              "role_name_or_id",
	PermittedByDefault: true,
	ParamsAllowed:      true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Make sure request is not empty.
		if q == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "", "")
			return err
		}

		// Try to get guild of question.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			return err
		}

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, m, q)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, m, roles, err)
		}

		// Try to assign role.
		err = sg.GuildMemberRoleAdd(guild.ID, m.Author.ID, roles[0].ID)
		if err != nil {
			_, err = sg.RespondDanger(m, "", err.Error())
			return err
		}

		// Respond about successful operation.
		_, err = sg.RespondSuccess(m, "", "you now have `"+roles[0].Name+"` role")
		return err
	},
}
