package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var leaveCmd = &sugo.Command{
	Trigger:            "leave",
	Description:        "Removes person to the public role.",
	Usage:              "rolenameorid",
	PermittedByDefault: true,
	ParamsAllowed: true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Make sure request is not empty.
		if q == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return err
		}

		// Try to get guild of question.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			return err
		}

		// Try to find user public role based on query.
		roles, err := publicRoles.findUserPublicRole(sg, m, q)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, m, roles, err)
		}

		// Try to remove user role.
		err = sg.GuildMemberRoleRemove(guild.ID, m.Author.ID, roles[0].ID)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return err
		}

		// Respond about operation being successful.
		_, err = sg.RespondSuccessMention(m, "you don't have `"+roles[0].Name+"` role any more")
		return err
	},
}
