package public_roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"sort"
)

var publicRoles = tPublicRoles{}

func respondFuzzyRolesSearchIssue(sg *sugo.Instance, m *discordgo.Message, roles []*discordgo.Role, err error) error {
	// Start building response.
	var response string
	response = err.Error()

	// If we have got at least one suggested role.
	if len(roles) > 0 {
		// Make an array of suggested role names.
		var suggestedRoles []*discordgo.Role
		for _, role := range roles {
			suggestedRoles = append(suggestedRoles, role)
		}
		response = response + ", try these:\n```\n"
		response = response + sugo.FmtStringsSlice(rolesToRoleNames(suggestedRoles), ", ", 1500, "\n...", "")
		response = response + "```"
	}

	_, err = sg.RespondWarning(m, "", response)
	return err
}

func rolesToRoleNames(roles []*discordgo.Role) []string {
	var roleNames = make([]string, 0)
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	sort.Strings(roleNames)
	return roleNames
}

// Module allows to manipulate public roles.
var Module = &sugo.Module{
	Startup:     startup,
	RootCommand: rootCommand,
}
