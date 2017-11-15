package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

var addCmd = &sugo.Command{
	Trigger:     "add",
	Description: "Makes existing role public.",
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

		// Get all guild roles.
		roles, err := sg.GuildRoles(guild.ID)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return
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
					_, err = sg.RespondFailMention(
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
			_, err = sg.RespondFailMention(m, "no roles found for query")
			return
		}

		// Otherwise add new role to the public roles list.
		err = storage.addGuildPublicRole(guild.ID, matchedRole)
		if err != nil {
			_, err = sg.RespondFailMention(m, err.Error())
			return
		}

		// And notify user about success.
		_, err = sg.RespondSuccessMention(m, "role `"+matchedRole.Name+"` is public now")
		return
	},
}
