package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var createCmd = &sugo.Command{
	Trigger:     "create",
	Description: "Creates new role on server and makes it public.",
	Usage:       "rolename",
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

		// Get all guild roles.
		roles, err := sg.GuildRoles(guild.ID)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return
		}

		// Try to match role.
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(q) {
				// We have found the role with the same name.
				_, err = sg.RespondFailMention(
					m,
					"role with such name already exists",
				)
				return err
			}
		}

		// If we did not find any match, try to create new role.
		role, err := sg.GuildRoleCreate(guild.ID)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return
		}

		// Set new role properties.
		sg.GuildRoleEdit(guild.ID, role.ID, q, 0, false, 0, true)

		// And add new role to the list of the public roles.
		storage.addGuildPublicRole(guild.ID, role)

		// And notify user about success.
		_, err = sg.RespondSuccessMention(m, "new role `"+role.Name+"` was created and made public")
		return
	},
}
