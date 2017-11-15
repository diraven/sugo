package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var delCmd = &sugo.Command{
	Trigger:     "del",
	Description: "Makes given role not public (does not delete the role itself).",
	Usage:       "rolenameorid",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure query is not empty.
		if q == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		// Sync roles with the server.
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

		// Convert one-item-map to roleID and roleName
		roleID, roleName := pickRoleFromMap(roles)

		// Delete role.
		storage.delGuildPublicRole(guild.ID, roleID)

		// Notify user about success of the operation.
		_, err = sg.RespondSuccessMention(m, "role `"+roleName+"` is not public any more")
		return
	},
}
