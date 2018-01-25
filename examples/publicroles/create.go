package publicroles

import (
	"github.com/diraven/sugo"
	"strings"
	"github.com/bwmarrin/discordgo"
)

var createCmd = &sugo.Command{
	Trigger:     "create",
	Description: "Creates new role on server and makes it public.",
	PermissionsRequired: discordgo.PermissionManageRoles,
	HasParams:           true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		// Make sure query is not empty.
		if req.Query == "" {
			_, err = sg.RespondBadCommandUsage(req, "", "")
			return err
		}

		// Get guild.
		guild, err := req.GetGuild()
		if err != nil {
			return err
		}

		// Get all guild roles.
		roles, err := sg.Session.GuildRoles(guild.ID)
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// Try to match role.
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(req.Query) {
				// We have found the role with the same name.
				_, err = sg.RespondDanger(
					req, "",
					"role with such name already exists",
				)
				return err
			}
		}

		// If we did not find any match, try to create new role.
		role, err := sg.Session.GuildRoleCreate(guild.ID)
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// Set new role properties.
		role, err = sg.Session.GuildRoleEdit(guild.ID, role.ID, req.Query, 0, false, 0, true)
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// And add new role to the list of the public roles.
		// Otherwise add new role to the public roles list.
		err = storage.add(role.ID)
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}

		// Save our changes.
		err = storage.save()
		if err != nil {
			_, err = sg.RespondDanger(req, "", err.Error())
			return err
		}


		// And notify user about success.
		_, err = sg.RespondSuccess(req, "", "new role `"+role.Name+"` was created and made public")
		return err
	},
}
