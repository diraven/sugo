package public_roles

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"io/ioutil"
	"os"
	"sort"
)

var storage *sStorage

var DATA_FILENAME = "public_roles_v2.json"

func init() {
	storage = &sStorage{
		make(map[string][]string),
	}
}

func respondFuzzyRolesSearchIssue(sg *sugo.Instance, m *discordgo.Message, roles []*discordgo.Role, err error) error {
	// Start building response.
	var response string
	response = err.Error()

	// If we have got at least one suggested role.
	if len(roles) > 0 {
		// Make an array of suggested role names.
		suggestedRoles := []*discordgo.Role{}
		for _, role := range roles {
			suggestedRoles = append(suggestedRoles, role)
		}
		response = response + ", try these:\n```\n"
		response = response + sugo.FmtStringsSlice(rolesToRoleNames(suggestedRoles), ", ", 1500, "\n...", "")
		response = response + "```"
	}

	_, err = sg.RespondFailMention(m, response)
	return err
}

func rolesToRoleNames(roles []*discordgo.Role) []string {
	var roleNames []string = []string{}
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	sort.Strings(roleNames)
	return roleNames
}

// CmdRSS allows to manipulate public roles.
var Cmd = &sugo.Command{
	Trigger:            "pr",
	Description:        "Allows to manipulate public roles.",
	PermittedByDefault: true,
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Try to find role based on query.
		roles, err := storage.findGuildPublicRole(sg, m, q)

		// Start building response.
		var response string

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response = response + "```\n"
			response = response + sugo.FmtStringsSlice(rolesToRoleNames(roles), "\n", 1500, "\n...", "")
			response = response + "```"
			_, err = sg.RespondTextMention(m, response)
		} else {
			_, err = sg.RespondTextMention(m, "nothing found")
		}

		return err
	},
	SubCommands: []*sugo.Command{
		myCmd,
		whoCmd,
		addCmd,
		delCmd,
		joinCmd,
		leaveCmd,
		createCmd,
		statsCmd,
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
