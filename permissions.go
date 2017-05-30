package sugo

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type iPermissionsStorage interface {
	get(sg *Instance, commandPath string, roleID string) (isAllowed bool, exists bool)
	set(sg *Instance, commandPath string, roleID string, isAllowed bool)
	default_(sg *Instance, commandPath string, roleID string)
	load(sg *Instance) (data_length int, err error)
	save(sg *Instance) (data_length int, err error)
	startup(sg *Instance) error
	teardown(sg *Instance) error
}

const PERMISSIONS_DATA_FILENAME = "permissions.json"

type permissionStorage struct {
	Permissions map[string]bool
}

func (p *permissionStorage) get(sg *Instance, commandPath string, roleID string) (isAllowed bool, exists bool) {
	var key string
	key = roleID + ":" + commandPath
	isAllowed, ok := p.Permissions[key]
	if ok {
		return isAllowed, true
	} else {
		return
	}
}

func (p *permissionStorage) set(sg *Instance, commandPath string, roleID string, isAllowed bool) {
	var key string
	key = roleID + ":" + commandPath
	p.Permissions[key] = isAllowed
}

func (p *permissionStorage) default_(sg *Instance, commandPath string, roleID string) {
	var key string
	key = roleID + ":" + commandPath
	delete(p.Permissions, key)
}

func (p *permissionStorage) load(sg *Instance) (data_length int, err error) {
	if _, error_type := os.Stat(PERMISSIONS_DATA_FILENAME); os.IsNotExist(error_type) {
		sg.DebugLog(0, "No perms file found. Empty storage initialized.")
		// File to load data from does not exist.
		// Check if perms storage is empty and initialize it.
		permsStorage := sg.permissions.(*permissionStorage)
		if permsStorage.Permissions == nil {
			permsStorage.Permissions = make(map[string]bool)
		}
		return
	}

	// Load file.
	data, err := ioutil.ReadFile(PERMISSIONS_DATA_FILENAME)
	if err != nil {
		return
	}

	// Decode JSON data.
	json.Unmarshal(data, sg.permissions.(*permissionStorage))
	if err != nil {
		return
	}

	// Log the operation results.
	data_length = len(data)
	sg.DebugLog(0, "Permissions loaded successfully,", data_length, "bytes read.")

	return
}

func (p *permissionStorage) save(sg *Instance) (data_length int, err error) {
	// Encode our data into JSON.
	data, err := json.Marshal(sg.permissions.(*permissionStorage))
	if err != nil {
		return
	}

	// Save data into file.
	err = ioutil.WriteFile(PERMISSIONS_DATA_FILENAME, data, 0644)
	if err != nil {
		return
	}

	data_length = len(data)
	sg.DebugLog(0, "Permissions saved successfully,", data_length, "bytes written.")

	return
}

func (p *permissionStorage) startup(sg *Instance) (err error) {
	_, err = p.load(sg)
	return
}

func (p *permissionStorage) teardown(sg *Instance) (err error) {
	_, err = p.save(sg)
	return
}

func findRole(sg *Instance, m *discordgo.Message, oldQ string) (role *discordgo.Role, q string) {
	// Extract role ID from the query string.
	ss := strings.Split(oldQ, " ")
	roleMention := ss[0]
	roleID := strings.TrimLeft(roleMention, "<#@&")
	roleID = strings.TrimRight(roleID, ">")
	q = strings.Replace(oldQ, roleMention, "", 1)
	q = strings.TrimSpace(q)

	// Get channel from the message.
	channel, err := sg.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Get guild from the channel.
	g, err := sg.State.Guild(channel.GuildID)
	if err != nil {
		return
	}

	// Try to find role specified.
	for _, r := range g.Roles {
		if r.ID == roleID {
			role = r
			break
		}
	}
	return
}

// Help shows help section for appropriate command.
var CmdPerms = &Command{
	Trigger:     "perms",
	RootOnly:    true,
	Description: "Allows to manipulate custom command permissions.",
	SubCommands: []*Command{
		{
			Trigger:     "load",
			Description: "Allows specific command usage for specific role.",
			Usage:       "@role command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				data_length, err := sg.permissions.load(sg)
				if err != nil {
					return
				}
				sg.RespondTextMention(m, "Permissions loaded. "+strconv.FormatInt(int64(data_length), 10)+" bytes read.")
				return
			},
		},
		{
			Trigger:     "save",
			Description: "Allows specific command usage for specific role.",
			Usage:       "@role command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				data_length, err := sg.permissions.save(sg)
				if err != nil {
					return
				}
				sg.RespondTextMention(m, "Permissions saved. "+strconv.FormatInt(int64(data_length), 10)+" bytes written.")
				return
			},
		},
		{
			Trigger:     "allow",
			Description: "Allows specific command usage for specific role.",
			Usage:       "@role command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				// Try to find role.
				role, q := findRole(sg, m, q)
				if role == nil {
					sg.RespondTextMention(m, "Role not found.")
					return
				}

				// Try to find command.
				command, err := sg.rootCommand.search(sg, q, m)
				if command == nil {
					sg.respondCommandNotFound(m)
					return
				}

				sg.permissions.set(sg, command.path(), role.ID, true)
				sg.RespondTextMention(m, "`"+role.Name+"` is now allowed to use \""+q+"\".")
				return
			},
		},
		{
			Trigger:     "deny",
			Description: "Denies specific command usage for specific role.",
			Usage:       "@role command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				// Try to find role.
				role, q := findRole(sg, m, q)
				if role == nil {
					sg.RespondTextMention(m, "Role not found.")
					return
				}

				// Try to find command.
				command, err := sg.rootCommand.search(sg, q, m)
				if command == nil {
					sg.respondCommandNotFound(m)
					return
				}

				sg.permissions.set(sg, command.path(), role.ID, false)
				sg.RespondTextMention(m, "`"+role.Name+"` is now not allowed to use \""+q+"\".")
				return
			},
		},
		{
			Trigger:     "default",
			Description: "Sets permissions for specific role usage by role specified to default.",
			Usage:       "@role command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				// Try to find role.
				role, q := findRole(sg, m, q)
				if role == nil {
					sg.RespondTextMention(m, "Role not found.")
					return
				}

				// Try to find command.
				command, err := sg.rootCommand.search(sg, q, m)
				if command == nil {
					sg.respondCommandNotFound(m)
					return
				}

				sg.permissions.default_(sg, command.path(), role.ID)
				sg.RespondTextMention(m, "`"+q+"` usage permissions by `"+role.Name+"` were set to default.")
				return
			},
		},
	},
}
