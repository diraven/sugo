package public_roles

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"strings"
)

type sStorage struct {
	RoleIDs map[string][]string // Items[GuildID][]RoleID
}

func (s *sStorage) getGuildPublicRoles(sg *sugo.Instance, m *discordgo.Message) discordgo.Roles {
	// Make storage to store all public roles we discovered user is in.
	roles := discordgo.Roles{}

	// Get guild.
	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return roles
	}

	// Get all guild roles.
	guildRoles, err := sg.GuildRoles(guild.ID)
	if err != nil {
		return roles
	}

	// Check if guild exists and has public roles.
	roleIDs, ok := s.RoleIDs[guild.ID]

	// If guild exists - try to match our saved roles with actual guild roles.
	if ok {
		// Variable to hold roles that were deleted on the server.
		var removedRolesIDs = []string{}

		// Variable to hold a switch if we found a role or not.
		var roleFound bool

		// For each roleID we have stored.
		for _, roleID := range roleIDs {
			roleFound = false
			// Try to get it from the guild roles.
			for _, guildRole := range guildRoles {
				if roleID == guildRole.ID {
					roles = append(roles, guildRole)
					roleFound = true
					break
				}
			}
			if !roleFound {
				// We did not find the role in the server roles list, which means the role was deleted on the server side.
				removedRolesIDs = append(removedRolesIDs, roleID)
			}
		}

		// Clean up our role storage to remove references to the roles that do not exist any more.
		for _, removedRoleID := range removedRolesIDs {
			s.delGuildPublicRole(guild.ID, removedRoleID)
		}
	}

	return roles
}

func (s *sStorage) getUserPublicRoles(sg *sugo.Instance, m *discordgo.Message) discordgo.Roles {
	// Make storage to store all public roles we discovered user is in.
	roles := discordgo.Roles{}

	// Get guild member.
	member, err := sg.MemberFromMessage(m)
	if err != nil {
		return roles
	}

	// Get guild public roles.
	guildPublicRoles := storage.getGuildPublicRoles(sg, m)

	// Iterate over all member roles.
	for _, memberRoleID := range member.Roles {
		// Iterate over guild public roles.
		for _, role := range guildPublicRoles {
			// If member role ID is in the guild public roles list.
			if memberRoleID == role.ID {
				// Then we have found a member public role and add it to the list.
				roles = append(roles, role)
				break
			}
		}
	}

	// Return resulting user public roles list.
	return roles
}

func (s *sStorage) addGuildPublicRole(guildID string, role *discordgo.Role) (err error) {
	// Check if guild exists.
	_, ok := s.RoleIDs[guildID]
	if !ok {
		// if guild does not exist - add new one.
		s.RoleIDs[guildID] = []string{}
	}

	// Check if role is public.
	for _, roleID := range s.RoleIDs[guildID] {
		if role.ID == roleID {
			return errors.New("this role is already public")
		}
	}

	// Make role public.
	s.RoleIDs[guildID] = append(s.RoleIDs[guildID], role.ID)
	return
}

// Function uses roleID instead of *discordgo.Role, because role may already not be on the server when we try to
// delete it, so we won't be able to retrieve it's properties.
func (s *sStorage) delGuildPublicRole(guildID string, roleID string) (err error) {
	// Variable that stores key existence check results.
	var ok bool

	// Check if guild exists.
	_, ok = s.RoleIDs[guildID]
	if !ok {
		// if guild does not exist - we do nothing.
		return
	}

	// Check if role is public.
	var idx int = -1
	for i, storedRoleID := range s.RoleIDs[guildID] {
		if roleID == storedRoleID {
			// If we have found a public role to be deleted - save it's index.
			idx = i
		}
	}

	if idx >= 0 {
		// Now delete item with the given index.
		s.RoleIDs[guildID] = append(s.RoleIDs[guildID][:idx], s.RoleIDs[guildID][idx+1:]...)
		return
	}

	// If guild exists and role does not exist - return error.
	return errors.New("role not found")
}

func (s *sStorage) findRole(roles discordgo.Roles, q string) (suggestedRoles discordgo.Roles, err error) {
	// Edit distance considered similar enough.
	var expectedEditDistance int = 2

	// Initialize suggested roles slice.
	suggestedRoles = discordgo.Roles{}

	// Iterate over stored roles.
	for _, role := range roles {
		// If we have detected a perfect fit for the query:
		if q == role.ID || strings.Contains(strings.ToLower(role.Name), strings.ToLower(q)) {
			// Add the role we found into suggestions.
			suggestedRoles = append(suggestedRoles, role)
		}
	}

	// If amount of suggested role IDs is exactly one - we have got a perfect fit.
	if len(suggestedRoles) == 1 {
		return suggestedRoles, nil
	}

	// If amount of suggested role IDs is more then one - we have found more then 1 role fitting the query.
	if len(suggestedRoles) > 1 {
		return suggestedRoles, errors.New("multiple roles found")
	}

	// If amount of suggested role IDs is 0 - we have found no fitting roles.
	if len(suggestedRoles) == 0 {
		// Try to figure out what did the user want by calculating Levenshtein (edit) distance between query and role names.
		for _, role := range roles {
			// Variable to hold edit distance.
			var d int

			if strings.Contains(role.Name, " ") && !strings.Contains(q, " ") {
				// If role Name is multi word while query is not, we will try to match query with every word of role name.
				for _, roleNameWord := range strings.Split(role.Name, " ") {
					d = levenshtein.DistanceForStrings([]rune(strings.ToLower(roleNameWord)), []rune(strings.ToLower(q)), levenshtein.DefaultOptions)
					// If edit distance is less then equal then expected:
					if d <= expectedEditDistance {
						// Add the role id to the suggested list.
						suggestedRoles = append(suggestedRoles, role)
						break
					}
				}
			} else {
				// Otherwise just try to match full query with full role name.
				// Calculate edit distance between full query and full role name.
				d = levenshtein.DistanceForStrings([]rune(strings.ToLower(role.Name)), []rune(strings.ToLower(q)), levenshtein.DefaultOptions)
				// If edit distance is small enough, we consider role to be a suggestion.
				if d <= expectedEditDistance {
					suggestedRoles = append(suggestedRoles, role)
				}
			}
		}

		// Now return what we got, even if we have found no suggestions.
		return suggestedRoles, errors.New("nothing found")
	}

	// This should never happen as map length can not be negative.
	panic(suggestedRoles)
}

func (s *sStorage) findUserPublicRole(sg *sugo.Instance, m *discordgo.Message, q string) (roles discordgo.Roles, err error) {
	return s.findRole(s.getUserPublicRoles(sg, m), q)
}

func (s *sStorage) findGuildPublicRole(sg *sugo.Instance, m *discordgo.Message, q string) (roles discordgo.Roles, err error) {
	return s.findRole(s.getGuildPublicRoles(sg, m), q)
}
