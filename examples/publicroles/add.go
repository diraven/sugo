package publicroles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var addCmd = &sugo.Command{
	Trigger:             "add",
	Description:         "Makes existing role public.",
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

		// Process request.
		var request string

		if len(req.Message.MentionRoles) > 0 {
			// If there is at least one role mention - we use that mention.
			request = req.Message.MentionRoles[0]
		} else {
			// Otherwise we just take full request.
			request = req.Query
		}

		// Make a storage for role we matched.
		var matchedRole *discordgo.Role

		// Try to match role.
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(request) || role.ID == request {
				if matchedRole != nil {
					_, err = sg.RespondDanger(
						req, "", "too many roles found, try again with a different search",
					)
					return err
				}
				matchedRole = role
			}
		}

		// If we did not find any match:
		if matchedRole == nil {
			// Notify user about fail.
			_, err = sg.RespondDanger(req, "", "no roles found for query")
			return err
		}

		// Otherwise add new role to the public roles list.
		err = storage.add(matchedRole.ID)
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
		_, err = sg.RespondSuccess(req, "", "role `"+matchedRole.Name+"` is public now")
		return err
	},
}
