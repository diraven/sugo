package public_roles

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
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

func (s *sStorage) getGuildPublicRoles(guildID string) (map[string]string) {
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

func (s *sStorage) findGuildPublicRole(sg *sugo.Instance, m *discordgo.Message, q string) (roleID string, err error) {
	// Get a guild.
	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return
	}

	// Prepare role storage.
	var foundRoleID string
	// Iterate over stored roles.
	for storedID, storedName := range s.getGuildPublicRoles(guild.ID) {
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