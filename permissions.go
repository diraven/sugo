package sugo

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type iPermissionStorage interface {
	get(sg *Instance, commandPath string, roleID string) (isAllowed bool, exists bool)
	set(sg *Instance, commandPath string, roleID string, isAllowed bool)
	default_(sg *Instance, commandPath string, roleID string)
	load(sg *Instance) (data_length int, err error)
	save(sg *Instance) (data_length int, err error)
	startup(sg *Instance) error
	teardown(sg *Instance) error
}

const DATA_FILENAME = "permissions.json"

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
	if _, error_type := os.Stat(DATA_FILENAME); os.IsNotExist(error_type) {
		log.Println("No perms file found. Empty storage initialized.")
		// File to load data from does not exist.
		// Check if perms storage is empty and initialize it.
		permsStorage := sg.permissions.(*permissionStorage)
		if permsStorage.Permissions == nil {
			permsStorage.Permissions = make(map[string]bool)
		}
		return
	}

	// Load file.
	data, err := ioutil.ReadFile(DATA_FILENAME)
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
	log.Println("Permissions loaded successfully,", data_length, "bytes read.")

	return
}

func (p *permissionStorage) save(sg *Instance) (data_length int, err error) {
	// Encode our data into JSON.
	data, err := json.Marshal(sg.permissions.(*permissionStorage))
	if err != nil {
		return
	}

	// Save data into file.
	err = ioutil.WriteFile(DATA_FILENAME, data, 0644)
	if err != nil {
		return
	}

	data_length = len(data)
	log.Println("Permissions saved successfully,", data_length, "bytes written.")

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
				command, err := findCommand(q, m, Bot.rootCommand.SubCommands)
				if command == nil {
					sg.RespondTextMention(m, "Command not found.")
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
				command, err := findCommand(q, m, Bot.rootCommand.SubCommands)
				if command == nil {
					sg.RespondTextMention(m, "Command not found.")
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
				command, err := findCommand(q, m, Bot.rootCommand.SubCommands)
				if command == nil {
					sg.RespondTextMention(m, "Command not found.")
					return
				}

				sg.permissions.default_(sg, command.path(), role.ID)
				sg.RespondTextMention(m, "`"+q+"` usage permissions by `"+role.Name+"` were set to default.")
				return
			},
		},
		{
			Trigger:     "debug",
			Description: "Prints out permission check process for the given user and command.",
			Usage:       "@user command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				// Get user.
				ss := strings.Split(q, " ")
				userID := ss[0]
				q = strings.TrimSpace(strings.Replace(q, userID, "", 1))
				userID = strings.TrimLeft(userID, "<@!")
				userID = strings.TrimRight(userID, ">")

				// Get command.
				command, err := findCommand(q, m, sg.rootCommand.SubCommands)
				if err != nil {
					return
				}
				if command == nil {
					sg.RespondTextMention(m, "Command not found.")
					return
				}

				// Get channel.
				channel, err := sg.State.Channel(m.ChannelID)
				if err != nil {
					return
				}

				// Get guild.
				guild, err := sg.State.Guild(channel.GuildID)
				if err != nil {
					return
				}

				// Get guild member.
				member, err := sg.State.Member(guild.ID, userID)
				if err != nil {
					sg.RespondTextMention(m, "User not found.")
					return
				}

				// Start building response embed.
				embed := &discordgo.MessageEmbed{
					Title: "Permissions for: " + member.User.Username,
					//Description: strings.Join(sg.triggers(), ", "),
					Color:  ColorInfo,
					Fields: []*discordgo.MessageEmbedField{},
				}

				// Get all member roles and check if any of those have custom permissions.
				var isAllowed bool
				var exists bool
				var role *discordgo.Role
				for _, r := range member.Roles {
					isAllowed, exists = sg.permissions.get(sg, q, r)

					// Check everyone role.
					isAllowed, exists = sg.permissions.get(sg, q, guild.ID)
					if exists {
						embed.Fields = append(embed.Fields,
							&discordgo.MessageEmbedField{
								Name:   "@everyone:",
								Value:  strconv.FormatBool(isAllowed),
								Inline: true,
							},
						)
					}
					// Check default.
					embed.Fields = append(embed.Fields,
						&discordgo.MessageEmbedField{
							Name:   "Default:",
							Value:  strconv.FormatBool(command.PermittedByDefault),
							Inline: true,
						},
					)
					// Check if user is root.
					embed.Fields = append(embed.Fields,
						&discordgo.MessageEmbedField{
							Name:   "Is Root:",
							Value:  strconv.FormatBool(sg.isRoot(member.User)),
							Inline: true,
						},
					)

					if exists {
						role, err = sg.State.Role(guild.ID, r)
						if err != nil {
							return
						}
						embed.Fields = append(embed.Fields,
							&discordgo.MessageEmbedField{
								Name:   "`" + role.Name + "` (" + strconv.FormatInt(int64(role.Position), 10) + ")",
								Value:  strconv.FormatBool(isAllowed),
								Inline: true,
							},
						)
					}
				}
				sg.RespondEmbed(m, embed)
				return
			},
		},
	},
}
