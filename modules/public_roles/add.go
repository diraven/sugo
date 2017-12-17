package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var addCmd = &sugo.Command{
	Trigger:       "add",
	Description:   "Makes existing role public.",
	Usage:         "role_name_or_id",
	ParamsAllowed: true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Make sure query is not empty.
		if q == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return err
		}

		// Get a guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			_, err = sg.RespondDanger(m, err.Error())
			return err
		}

		// Get all guild roles.
		roles, err := sg.GuildRoles(guild.ID)
		if err != nil {
			_, err = sg.RespondDanger(m, err.Error())
			return err
		}

		// Process request.
		var request string

		if len(m.MentionRoles) > 0 {
			// If there is at least one role mention - we use that mention.
			request = m.MentionRoles[0]
		} else {
			// Otherwise we just take full request.
			request = q
		}

		// Make a storage for role we matched.
		var matchedRole *discordgo.Role

		// Try to match role.
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(request) || role.ID == request {
				if matchedRole != nil {
					_, err = sg.RespondDanger(
						m,
						"too many roles found, try again with a different search",
					)
					return err
				} else {
					matchedRole = role
				}
			}
		}

		// If we did not find any match:
		if matchedRole == nil {
			// Notify user about fail.
			_, err = sg.RespondDanger(m, "no roles found for query")
			return err
		}

		// Otherwise add new role to the public roles list.
		err = publicRoles.add(sg, guild.ID, matchedRole.ID)
		if err != nil {
			_, err = sg.RespondDanger(m, err.Error())
			return err
		}

		// And notify user about success.
		_, err = sg.RespondSuccess(m, "role `"+matchedRole.Name+"` is public now")
		return err
	},
}
