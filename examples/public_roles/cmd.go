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

		// Try to find role based on query.
		_, suggestedRoleIDs, _ := storage.findGuildPublicRole(sg, m, q)

		// Start building response.
		var response string

		// If we have got at least one suggested role.
		if len(suggestedRoleIDs) > 0 {
			// Make an array of suggested role names.
			suggestedRoleNames := []string{}
			for _, id := range suggestedRoleIDs {
				name, err := storage.getPublicRoleName(guild.ID, id)
				if err != nil {
					return err
				}
				suggestedRoleNames = append(suggestedRoleNames, name)
			}

			response = response + "here is what I found based on your query\n```"
			response = response + sugo.FmtStringsSlice(suggestedRoleNames, "\n", 1500, "\n...", "")
			response = response + "```"
			_, err = sg.RespondTextMention(m, response)
		} else {
			_, err = sg.RespondTextMention(m, "no public roles found")
		}

		return err
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
				storedRoleID, suggestedRoleIDs, err := storage.findGuildPublicRole(sg, m, q)
				if err != nil {
					// Start building response.
					var response string
					response = err.Error()

					// If we have got at least one suggested role.
					if len(suggestedRoleIDs) > 0 {
						// Make an array of suggested role names.
						suggestedRoleNames := []string{}
						for _, id := range suggestedRoleIDs {
							name, err := storage.getPublicRoleName(guild.ID, id)
							if err != nil {
								return err
							}
							suggestedRoleNames = append(suggestedRoleNames, name)
						}
						response = response + "\n Did you mean one of following:\n```"
						response = response + sugo.FmtStringsSlice(suggestedRoleNames, ", ", 1500, "\n...", "")
						response = response + "```"
					}

					_, err = sg.RespondFailMention(m, response)
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
				storage.addGuildPublicRole(guild.ID, matchedRole.ID, matchedRole.Name)

				// And notify user about success.
				_, err = sg.RespondSuccessMention(m, "role `"+matchedRole.Name+"` is public now")
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
				roleID, suggestedRoleIDs, err := storage.findGuildPublicRole(sg, m, q)
				if err != nil {
					// Start building response.
					var response string
					response = err.Error()

					// If we have got at least one suggested role.
					if len(suggestedRoleIDs) > 0 {
						// Make an array of suggested role names.
						suggestedRoleNames := []string{}
						for _, id := range suggestedRoleIDs {
							name, err := storage.getPublicRoleName(guild.ID, id)
							if err != nil {
								return err
							}
							suggestedRoleNames = append(suggestedRoleNames, name)
						}
						response = response + "\n Did you mean one of following:\n```"
						response = response + sugo.FmtStringsSlice(suggestedRoleNames, ", ", 1500, "\n...", "")
						response = response + "```"
					}

					_, err = sg.RespondFailMention(m, response)
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
				_, err = sg.RespondSuccessMention(m, "role `"+roleName+"` is not public any more")
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

				// Try to get guild of question.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return
				}

				// Try to find role based on query.
				storedRoleID, suggestedRoleIDs, err := storage.findGuildPublicRole(sg, m, q)
				if err != nil {
					// Start building response.
					var response string
					response = err.Error()

					// If we have got at least one suggested role.
					if len(suggestedRoleIDs) > 0 {
						// Make an array of suggested role names.
						suggestedRoleNames := []string{}
						for _, id := range suggestedRoleIDs {
							name, err := storage.getPublicRoleName(guild.ID, id)
							if err != nil {
								return err
							}
							suggestedRoleNames = append(suggestedRoleNames, name)
						}
						response = response + "\n Did you mean one of following:\n```"
						response = response + sugo.FmtStringsSlice(suggestedRoleNames, ", ", 1500, "\n...", "")
						response = response + "```"
					}

					_, err = sg.RespondFailMention(m, response)
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

				_, err = sg.RespondSuccessMention(m, "you now have `"+roleName+"` role")
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

				_, err = sg.RespondSuccessMention(m, "you don't have `"+roleName+"` role any more")
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
