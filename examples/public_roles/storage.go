package public_roles

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"sort"
	"strings"
)

type sStorage struct {
	Roles map[string]map[string]string // Items[GuildID][RoleID]=RoleName
}

func (s *sStorage) setPublicRoleName(guildID string, ruleID string, name string) (err error) {
	// Switch to store key existence check results.
	var ok bool

	// Check if given guild exists.
	_, ok = s.Roles[guildID]
	if !ok {
		return errors.New("guild not found")
	}

	// Check if given role exists within given guild.
	_, ok = s.Roles[guildID][ruleID]
	if !ok {
		return errors.New("rule not found")
	}

	// Assign role new name.
	s.Roles[guildID][ruleID] = name
	return
}

func (s *sStorage) getPublicRoleName(guildID string, ruleID string) (string, error) {
	// Switch to store key existence check results.
	var ok bool

	// Check if given guild exists.
	_, ok = s.Roles[guildID]
	if !ok {
		return "", errors.New("guild not found")
	}

	// Check if given role exists within given guild.
	_, ok = s.Roles[guildID][ruleID]
	if !ok {
		return "", errors.New("rule not found")
	}

	return s.Roles[guildID][ruleID], nil
}

func (s *sStorage) getGuildPublicRoles(guildID string) map[string]string {
	roles, ok := s.Roles[guildID]
	if ok {
		return roles
	} else {
		return make(map[string]string)
	}
}

func (s *sStorage) getUserPublicRoles(sg *sugo.Instance, m *discordgo.Message) (roles map[string]string, err error) {
	// Get guild member.
	member, err := sg.MemberFromMessage(m)
	if err != nil {
		return
	}

	// Make storage to store all public roles we discovered user is in.
	userPublicRoles := make(map[string]string)

	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return
	}

	// Get guild public roles.
	guildPublicRoles := storage.getGuildPublicRoles(guild.ID)

	// For each member role:
	for _, roleID := range member.Roles {
		// Make sure role is marked as public.
		storedName, ok := guildPublicRoles[roleID]
		if ok {
			// If role is marked as public - add it to the list.
			userPublicRoles[roleID] = storedName
		}
	}

	// Return resulting user public roles list.
	return userPublicRoles, nil
}

func (s *sStorage) addGuildPublicRole(guildID string, roleID string, roleName string) (err error) {
	// Variable that stores key existence check results.
	var ok bool

	// Check if guild exists.
	_, ok = s.Roles[guildID]
	if !ok {
		// if guild does not exist - add new one.
		s.Roles[guildID] = make(map[string]string)
	}

	// Check if role public.
	_, ok = s.Roles[guildID][roleID]
	if ok {
		// If role already exists - return error.
		return errors.New("the role is already public")
	}

	// If guild exists and role is not public - make role public.
	s.Roles[guildID][roleID] = roleName
	return
}

func (s *sStorage) delGuildPublicRole(guildID string, roleID string) (err error) {
	// Variable that stores key existence check results.
	var ok bool

	// Check if guild exists.
	_, ok = s.Roles[guildID]
	if !ok {
		// if guild does not exist - we do nothing.
		return
	}

	// Check if role public.
	_, ok = s.Roles[guildID][roleID]
	if ok {
		// If role exists - delete it.
		delete(s.Roles[guildID], roleID)
		return
	}

	// If guild exists and role does not exist - return error.
	return errors.New("role not found for deletion")
}

func (s *sStorage) syncPublicRoles(sg *sugo.Instance, m *discordgo.Message) (err error) {
	// Get a guild.
	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return
	}

	// Get all guild roles.
	roles, err := sg.GuildRoles(guild.ID)
	if err != nil {
		return
	}

	// Get guild public roles.
	guildPublicRoles := s.getGuildPublicRoles(guild.ID)

	// Switch to use to determine if role found or not.
	var roleFound bool

	// For each stored role:
	for publicRoleID := range guildPublicRoles {
		// For each guild role:
		for _, role := range roles {
			// If stored role ID is the same as guild role id:
			if publicRoleID == role.ID {
				// Then it still exists.
				roleFound = true
				// Update rule name.
				err = s.setPublicRoleName(guild.ID, role.ID, role.Name)
				if err != nil {
					return
				}
				break
			}
		}
		// If there is no matching ID of the stored role in the guild roles:
		if !roleFound {
			// We remove role from stored ones.
			s.delGuildPublicRole(guild.ID, publicRoleID)
		}
	}
	return
}

func (s *sStorage) findGuildPublicRole(sg *sugo.Instance, m *discordgo.Message, q string) (roleID string, suggestedRoleIDs []string, err error) {
	// Edit distance considered similar enough.
	var expectedEditDistance int = 4

	// Initialize suggested role ids slice.
	suggestedRoleIDs = []string{}

	// Get a guild.
	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return
	}

	// Iterate over stored roles.
	for storedID, storedName := range s.getGuildPublicRoles(guild.ID) {
		// If we have detected a perfect fit for the query:
		if q == storedID || strings.Contains(strings.ToLower(storedName), strings.ToLower(q)) {
			// Add the found role we found into suggestions.
			suggestedRoleIDs = append(suggestedRoleIDs, storedID)
		}
	}

	// If amount of suggested role IDs is exactly one - we have got a perfect fit.
	if len(suggestedRoleIDs) == 1 {
		return suggestedRoleIDs[0], suggestedRoleIDs, nil
	}

	// If amount of suggested role IDs is more then one - we have found more then 1 role fitting the query.
	if len(suggestedRoleIDs) > 1 {
		sort.Strings(suggestedRoleIDs)
		return "", suggestedRoleIDs, errors.New("multiple roles found for query")
	}

	// If amount of suggested role IDs is 0 - we have found no fitting roles.
	if len(suggestedRoleIDs) == 0 {
		// Try to figure out what did the user want by calculating Levenshtein (edit) distance between query and role names.
		for storedID, storedName := range s.getGuildPublicRoles(guild.ID) {
			// Variable to hold edit distance.
			var d int

			if strings.Contains(storedName, " ") && !strings.Contains(q, " ") {
				// If role Name is multi word while query is not, we will try to match query with every word of role name.
				for _, roleNameWord := range strings.Split(storedName, " ") {
					d = levenshtein.DistanceForStrings([]rune(strings.ToLower(roleNameWord)), []rune(strings.ToLower(q)), levenshtein.DefaultOptions)
					// If edit distance is less then equal then expected:
					if d <= expectedEditDistance {
						// Add the role id to the suggested list.
						suggestedRoleIDs = append(suggestedRoleIDs, storedID)
						break
					}
				}
			} else {
				// Otherwise just try to match full query with full role name.
				// Calculate edit distance between full query and full role name.
				d = levenshtein.DistanceForStrings([]rune(strings.ToLower(storedName)), []rune(strings.ToLower(q)), levenshtein.DefaultOptions)
				// If edit distance is small enough, we consider role to be a suggestion.
				if d <= expectedEditDistance {
					suggestedRoleIDs = append(suggestedRoleIDs, storedID)
				}
			}
		}

		// Now return what we got, even if we have found no suggestions.
		sort.Strings(suggestedRoleIDs)
		return "", suggestedRoleIDs, errors.New("no public roles found for query")
	}

	// This should never happen as slice length can not be negative.
	panic(suggestedRoleIDs)
}

func (s *sStorage) findUserPublicRole(sg *sugo.Instance, m *discordgo.Message, q string) (roleID string, err error) {
	// Prepare role storage.
	var foundRoleID string
	// Get all user public roles.
	roles, err := s.getUserPublicRoles(sg, m)
	// Iterate over stored roles.
	for storedID, storedName := range roles {
		if q == storedID || strings.Contains(strings.ToLower(storedName), strings.ToLower(q)) {
			if foundRoleID == "" {
				foundRoleID = storedID
			} else {
				return "", errors.New("multiple roles found for query")
			}
		}
	}
	if foundRoleID == "" {
		return "", errors.New("no public roles found for query")
	}
	return foundRoleID, nil
}
