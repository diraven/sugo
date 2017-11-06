package public_roles

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var storage = &sStorage{}

var DATA_FILENAME = "public_roles.json"

func init() {
	storage.Roles = make(map[string]map[string]string)
}

// CmdRSS allows to manipulate public roles.
var Cmd = &sugo.Command{
	Trigger:            "pr",
	Description:        "Allows to manipulate public roles.",
	PermittedByDefault: true,
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure our role list is in sync with the server.
		storage.syncPublicRoles(sg, m)

		// Get a guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			return
		}

		// Gather all role names.
		roleNames := []string{}
		for _, storedName := range storage.getGuildPublicRoles(guild.ID) {
			// Filter role names based on query if any.
			if q != "" {
				if strings.Contains(strings.ToLower(storedName), strings.ToLower(q)) {
					roleNames = append(roleNames, storedName)
				}
			} else {
				roleNames = append(roleNames, storedName)
			}
		}

		// Sort role names.
		sort.Strings(roleNames)

		// Start building response.
		var response string
		if q != "" {
			response = "Available Public Roles with \"" + q + "\" in name:\n```\n"
		} else {
			response = "Available Public Roles:\n```\n"
		}
		response = response + sugo.FmtStringsSlice(roleNames, "\n", 1500, "\n...", "")
		response = response + "```"

		_, err = sg.RespondTextMention(m, response)
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:            "my",
			Description:        "Lists public roles you are in.",
			PermittedByDefault: true,
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				// Make sure our role list is in sync with the server.
				storage.syncPublicRoles(sg, m)

				userPublicRoles, err := storage.getUserPublicRoles(sg, m)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Gather all role names.
				roleNames := []string{}
				for _, storedName := range userPublicRoles {
					// Filter role names based on query if any.
					if q != "" {
						if strings.Contains(strings.ToLower(storedName), strings.ToLower(q)) {
							roleNames = append(roleNames, storedName)
						}
					} else {
						roleNames = append(roleNames, storedName)
					}
				}

				// Sort role names.
				sort.Strings(roleNames)

				// Start building response.
				var response string
				if q != "" {
					response = "Your Public Roles with \"" + q + "\" in name:\n```\n"
				} else {
					response = "Your Public Roles:\n```\n"
				}
				response = response + sugo.FmtStringsSlice(roleNames, "\n", 1500, "\n...", "")
				response = response + "```"

				_, err = sg.RespondTextMention(m, response)
				return
			},
		},
		{
			Trigger:            "who",
			Description:        "Lists people that have public role specified.",
			PermittedByDefault: true,
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				// Make sure our role list is in sync with the server.
				storage.syncPublicRoles(sg, m)

				// Get a guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Try to find role based on query.
				storedRoleID, err := storage.findGuildPublicRole(sg, m, q)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Make members array we will be working with.
				memberNames := []string{}
				for _, member := range guild.Members {
					for _, roleID := range member.Roles {
						if roleID == storedRoleID {
							memberNames = append(memberNames, member.User.Username+"#"+member.User.Discriminator)
						}
					}
				}

				// Get stored role name.
				storedRoleName, err := storage.getPublicRoleName(guild.ID, storedRoleID)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Sort people.
				sort.Strings(memberNames)

				// Start building response.
				response := "People with public role \"" + storedRoleName + "\":\n```\n"
				response = response + sugo.FmtStringsSlice(memberNames, ", ", 1500, "...", ".")
				response = response + "```"

				_, err = sg.RespondTextMention(m, response)
				return
			},
		},
		{
			Trigger:     "add",
			Description: "Makes existing role public.",
			Usage:       "rolenameorid",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				// Make sure our role list is in sync with the server.
				storage.syncPublicRoles(sg, m)

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
					if strings.Contains(strings.ToLower(role.Name), request) || role.ID == request {
						if matchedRole != nil {
							_, err = sg.RespondFailMention(
								m,
								"Too many roles found... Try again with a different search.",
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
					_, err = sg.RespondFailMention(m, "no public roles found for query")
					return
				}

				// Otherwise add new role to the public roles list.
				storage.addGuildPublicRole(guild.ID, matchedRole.ID, matchedRole.Name)

				// And notify user about success.
				_, err = sg.RespondSuccessMention(m, "Role made public: "+matchedRole.Name+".")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Makes given role not public (does not delete the role itself).",
			Usage:       "rolenameorid",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				// Make sure query is not empty.
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				// Get a guild.
				// Get a guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Sync roles with the server.
				storage.syncPublicRoles(sg, m)

				// Try to find role based on query.
				roleID, err := storage.findGuildPublicRole(sg, m, q)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Get stored role name.
				roleName, err := storage.getPublicRoleName(guild.ID, roleID)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Delete role.
				storage.delGuildPublicRole(guild.ID, roleID)

				// Notify user about success of the operation.
				_, err = sg.RespondSuccessMention(m, "Role is not public any more: "+roleName+".")
				return
			},
		},
		{
			Trigger:            "join",
			Description:        "Joins person to the public role.",
			Usage:              "rolenameorid",
			PermittedByDefault: true,
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				// Make sure request is not empty.
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				// Sync public roles to make sure we won't try to assign role that does not exist any more.
				storage.syncPublicRoles(sg, m)

				// Try to find public role.
				storedRoleID, err := storage.findGuildPublicRole(sg, m, q)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Try to get guild of question.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return
				}

				// Try to assign role.
				err = sg.GuildMemberRoleAdd(guild.ID, m.Author.ID, storedRoleID)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Try to get role name.
				roleName, err := storage.getPublicRoleName(guild.ID, storedRoleID)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				_, err = sg.RespondSuccessMention(m, "You got new role \""+roleName+"\".")
				return
			},
		},
		{
			Trigger:            "leave",
			Description:        "Removes person to the public role.",
			Usage:              "rolenameorid",
			PermittedByDefault: true,
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				// Make sure request is not empty.
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				// Sync public roles to make sure we won't try to assign role that does not exist any more.
				storage.syncPublicRoles(sg, m)

				// Try to get guild of question.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return
				}

				// Try to find public role.
				storedRoleID, err := storage.findUserPublicRole(sg, m, q)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Try to remove user role.
				err = sg.GuildMemberRoleRemove(guild.ID, m.Author.ID, storedRoleID)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Try to get role name.
				roleName, err := storage.getPublicRoleName(guild.ID, storedRoleID)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				_, err = sg.RespondSuccessMention(m, "You have removed the role from yourself: "+roleName+".")
				return
			},
		},
	},
	Startup: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Check if file exists.
		if _, err = os.Stat(DATA_FILENAME); err == nil {
			// Load file.
			data, err := ioutil.ReadFile(DATA_FILENAME)
			if err != nil {
				return err
			}

			// Decode JSON data.
			json.Unmarshal(data, storage)
			if err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			// If there are any errors other then "file does not exist" - report error and shutdown.
			return
		}

		return nil

	},
	Teardown: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Encode our data into JSON.
		data, err := json.Marshal(storage)
		if err != nil {
			return
		}

		// Save data into file.
		err = ioutil.WriteFile(DATA_FILENAME, data, 0644)
		if err != nil {
			return
		}

		return
	},
}
