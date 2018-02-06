package publicroles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var createCmd = &sugo.Command{
	Trigger:             "create",
	Description:         "Creates new role on server and makes it public.",
	PermissionsRequired: discordgo.PermissionManageRoles,
	HasParams:           true,
	Execute: func(req *sugo.Request) error {
		var err error

		// Make sure query is not empty.
		if req.Query == "" {
			_, err = req.RespondBadCommandUsage("", "")
			return err
		}

		// Get guild.
		guild, err := req.GetGuild()
		if err != nil {
			return err
		}

		// Get all guild roles.
		roles, err := req.Sugo.Session.GuildRoles(guild.ID)
		if err != nil {
			_, err = req.RespondDanger("", err.Error())
			return err
		}

		// Try to match role.
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(req.Query) {
				// We have found the role with the same name.
				_, err = req.RespondDanger(
					"",
					"role with such name already exists",
				)
				return err
			}
		}

		// If we did not find any match, try to create new role.
		role, err := req.Sugo.Session.GuildRoleCreate(guild.ID)
		if err != nil {
			_, err = req.RespondDanger("", err.Error())
			return err
		}

		// Set new role properties.
		role, err = req.Sugo.Session.GuildRoleEdit(guild.ID, role.ID, req.Query, 0, false, 0, true)
		if err != nil {
			_, err = req.RespondDanger("", err.Error())
			return err
		}

		// And add new role to the list of the public roles.
		// Otherwise add new role to the public roles list.
		err = storage.add(role.ID)
		if err != nil {
			_, err = req.RespondDanger("", err.Error())
			return err
		}

		// Save our changes.
		err = storage.save()
		if err != nil {
			_, err = req.RespondDanger("", err.Error())
			return err
		}

		// And notify user about success.
		_, err = req.RespondSuccess("", "new role `"+role.Name+"` was created and made public")
		return err
	},
}
