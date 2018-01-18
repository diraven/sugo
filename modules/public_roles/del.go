package public_roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var delCmd = &sugo.Command{
	Trigger:     "del",
	Description: "Makes given role not public (does not delete the role itself).",
	Usage:       "role_name_or_id",
	AllowParams: true,
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Make sure query is not empty.
		if q == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "", "")
			return err
		}

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

		// Delete role.
		publicRoles.del(sg, guild.ID, roles[0].ID)

		// Notify user about success of the operation.
		_, err = sg.RespondSuccess(m, "", "role `"+roles[0].Name+"` is not public any more")
		return err
	},
}
