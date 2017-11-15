package public_roles

import (
	"context"
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

var joinCmd = &sugo.Command{
	Trigger:            "join",
	Description:        "Joins person to the public role.",
	Usage:              "rolenameorid",
	PermittedByDefault: true,
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure request is not empty.
		if q == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		// Sync public roles to make sure we won't try to assign role that does not exist any more.
		storage.syncPublicRoles(sg, m)

		// Try to get guild of question.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			return
		}

		// Try to find role based on query.
		roles, err := storage.findGuildPublicRole(sg, m, q)
		if err != nil {
			return respondFuzzyRolesSearchIssue(sg, m, roles, err)
		}

		// Convert one-item-map to roleID and roleName
		roleID, roleName := pickRoleFromMap(roles)

		// Try to assign role.
		err = sg.GuildMemberRoleAdd(guild.ID, m.Author.ID, roleID)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return
		}

		// Respond about successful operation.
		_, err = sg.RespondSuccessMention(m, "you now have `"+roleName+"` role")
		return
	},
}