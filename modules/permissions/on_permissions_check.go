package permissions

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

func onPermissionsCheck(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message) (*bool, error) {
	var err error
	var passed bool // the conclusion about if command is allowed

	// Get guild.
	guild, err := sg.GuildFromMessage(m)

	// Get guild member.
	member, err := sg.State.Member(guild.ID, m.Author.ID)
	if err != nil {
		passed = false
		return &passed, err // Just make sure we are safe and return false in case of any errors.
	}

	// Now we need to check if we have any custom settings for every role user has
	// sequentially starting from the topmost one.

	// For each role user has.
	var role *discordgo.Role
	var position = 0 // position of the custom role setting found
	var found bool   // if custom role setting found

	// Start with checking "everyone" role permissions.
	role, err = sg.State.Role(guild.ID, guild.ID)
	if err != nil {
		passed = false
		return &passed, err // Just make sure we are safe and return false in case of any errors.
	}

	isAllowed, exists := permissions.get(sg, c.Path(), role.ID)
	if exists {
		found = true
		passed = isAllowed
	}

	// And now check all the rest of the user roles.
	for _, roleID := range member.Roles {
		// Get role itself.
		role, err = sg.State.Role(guild.ID, roleID)
		if err != nil {
			passed = false
			return &passed, err // Just make sure we are safe and return false in case of any errors.
		}
		// Check if role is allowed to use the command.
		isAllowed, exists := permissions.get(sg, c.Path(), role.ID)

		// If custom role setting exists and it's position less then the one we have already found (role that is higher
		// takes precedence over the lower ones):
		if exists && role.Position >= position {
			position = role.Position // Store position of the role.
			found = true             // Make sure we know we have found a custom setting.
			passed = isAllowed       // Update return value.
		}
	}

	if found {
		// If we have found the custom role setting - we just return the one for the highest role we found.
		return &passed, nil
	}

	// We have found no setting for the given command, so we just return undefined result.
	return nil, nil
}