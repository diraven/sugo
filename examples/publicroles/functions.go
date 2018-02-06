package publicroles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
	"sort"
)

func respondFuzzyRolesSearchIssue(req *sugo.Request, roles []*discordgo.Role, err error) error {
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
		response = response + ", try these:\n\n"
		response = response + helpers.FmtStringsSlice(rolesToRoleNames(suggestedRoles), ", ", "`", 1500, "...", "")
	}

	_, err = req.RespondWarning("", response)
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
