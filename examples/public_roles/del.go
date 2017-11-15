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

		// Delete role.
		storage.delGuildPublicRole(guild.ID, roles[0].ID)

		// Notify user about success of the operation.
		_, err = sg.RespondSuccessMention(m, "role `"+roles[0].Name+"` is not public any more")
		return
	},
}
