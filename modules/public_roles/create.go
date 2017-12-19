package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var createCmd = &sugo.Command{
	Trigger:       "create",
	Description:   "Creates new role on server and makes it public.",
	Usage:         "role_name",
	ParamsAllowed: true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
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

		// Get all guild roles.
		roles, err := sg.GuildRoles(guild.ID)
		if err != nil {
			_, err = sg.RespondDanger(m, "", err.Error())
			return err
		}

		// Try to match role.
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(q) {
				// We have found the role with the same name.
				_, err = sg.RespondDanger(
					m, "",
					"role with such name already exists",
				)
				return err
			}
		}

		// If we did not find any match, try to create new role.
		role, err := sg.GuildRoleCreate(guild.ID)
		if err != nil {
			_, err = sg.RespondDanger(m, "", err.Error())
			return err
		}

		// Set new role properties.
		role, err = sg.GuildRoleEdit(guild.ID, role.ID, q, 0, false, 0, true)
		if err != nil {
			_, err = sg.RespondDanger(m, "", err.Error())
			return err
		}

		// And add new role to the list of the public roles.
		publicRoles.add(sg, guild.ID, role.ID)

		// And notify user about success.
		_, err = sg.RespondSuccess(m, "", "new role `"+role.Name+"` was created and made public")
		return err
	},
}
