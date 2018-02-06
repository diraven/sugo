package publicroles

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const storageFilename = "public_roles.yaml"

type publicRolesStorage []string // []roleID

var storage publicRolesStorage

func init() {
	storage = []string{}
}

// add makes role public.
func (pr *publicRolesStorage) add(roleID string) error {
	// Check if role is public.
	for _, storedRoleID := range *pr {
		if storedRoleID == roleID {
			return errors.New("this role is already public")
		}
	}

	// Append role.
	*pr = append(*pr, roleID)

	// Return no errors.
	return nil
}

// del makes rule not public.
func (pr *publicRolesStorage) del(roleID string) error {
	// Try to find role.
	storedRoleIDIdx := -1
	for i, storedRoleID := range *pr {
		if storedRoleID == roleID {
			storedRoleIDIdx = i
		}
	}

	// If role not found - return error.
	if storedRoleIDIdx < 0 {
		return errors.New("role not found")
	}

	// Remove role.
	store := *pr
	store[len(store)-1], store[storedRoleIDIdx] = store[storedRoleIDIdx], store[len(store)-1]
	*pr = store[:len(store)-1]

	// Return no errors.
	return nil
}

// reload reloads in-memory public_roles cache from the database.
func (pr *publicRolesStorage) save() error {
	var err error

	// Marshal storage.
	storageMarshalled, err := yaml.Marshal(pr)
	if err != nil {
		return err
	}

	// Save file.
	err = ioutil.WriteFile(storageFilename, storageMarshalled, 0644)
	if err != nil {
		return err
	}

	// Return no error.
	return nil
}

// reload reloads in-memory public_roles cache from the database.
func (pr *publicRolesStorage) load() error {
	// Variable to store errors if any.
	var err error

	// Initialize storage storage.
	*pr = publicRolesStorage{}

	// Open file.
	file, err := os.OpenFile(storageFilename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	// Load bytes from file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Unmarshal bytes into memory.
	err = yaml.Unmarshal(bytes, pr)
	if err != nil {
		return err
	}

	// Return no error.
	return nil
}

func (pr *publicRolesStorage) getGuildPublicRoles(req *sugo.Request) discordgo.Roles {
	// Make storage to store all roles we will find.
	roles := discordgo.Roles{}

	// Get guild.
	guild, err := req.GetGuild()
	if err != nil {
		return discordgo.Roles{}
	}

	// Get all guild roles.
	guildRoles, err := req.Sugo.Session.GuildRoles(guild.ID)
	if err != nil {
		return roles
	}

	// Variable to hold roles that were deleted on the server.
	var removedRolesIDs []string

	// Variable to hold a switch if we found a role or not.
	var roleFound bool

	// For each roleID we have stored.
	for _, roleID := range *pr {
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
		pr.del(removedRoleID)
	}

	// Save our changes.
	pr.save()

	return roles
}

func (pr *publicRolesStorage) getUserPublicRoles(req *sugo.Request) discordgo.Roles {
	// Make storage to store all public roles we discovered user is in.
	roles := discordgo.Roles{}

	// Get guild.
	guild, err := req.GetGuild()
	if err != nil {
		return discordgo.Roles{}
	}

	// Get guild member.
	member, err := req.Sugo.Session.State.Member(guild.ID, req.Message.Author.ID)
	if err != nil {
		return roles
	}

	// Get guild public roles.
	guildPublicRoles := pr.getGuildPublicRoles(req)

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

func (pr *publicRolesStorage) findRole(roles discordgo.Roles, q string) (discordgo.Roles, error) {
	// Edit distance considered similar enough.
	var expectedEditDistance = 2

	// Initialize suggested roles slice.
	suggestedRoles := discordgo.Roles{}

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

func (pr *publicRolesStorage) findUserPublicRole(req *sugo.Request, q string) (discordgo.Roles, error) {
	return pr.findRole(pr.getUserPublicRoles(req), q)
}

func (pr *publicRolesStorage) findGuildPublicRole(req *sugo.Request, q string) (discordgo.Roles, error) {
	return pr.findRole(pr.getGuildPublicRoles(req), q)
}
