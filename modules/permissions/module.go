package permissions

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"github.com/diraven/sugo"
)

var permissions = tPermissionsStorage{}

// findRole looks for a role in the query that also exists as a real role in given guild.
func findRole(sg *sugo.Instance, m *discordgo.Message, oldQ string) (role *discordgo.Role, q string) {
	// Extract role ID from the query string.
	ss := strings.Split(oldQ, " ")
	roleMention := ss[0]
	roleID := strings.TrimLeft(roleMention, "<#@&")
	roleID = strings.TrimRight(roleID, ">")
	q = strings.Replace(oldQ, roleMention, "", 1)
	q = strings.TrimSpace(q)

	// Get guild from the message.
	g, err := sg.GuildFromMessage(m)
	if err != nil {
		return nil, q
	}

	// Try to find role specified.
	for _, r := range g.Roles {
		if r.ID == roleID {
			role = r
			break
		}
	}
	return
}

// ModPerms handles custom permission checks set on the per guild basis. Called for every command in the chain to the
// bottommost one.
var Module = &sugo.Module{
	Startup:            startup,
	RootCommand:        rootCommand,
	OnPermissionsCheck: onPermissionsCheck,
}
