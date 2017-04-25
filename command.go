package sugo

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"github.com/diraven/sugo/helpers"
	"strings"
)

// Command struct describes basic command type.
type Command struct {
	// trigger is a sequence of symbols message should start with to match with the command.
	trigger string
	// rootOnly determines if the command is supposed to be used by root only.
	rootOnly bool
	// permissionsRequired is a slice of all permissions required by the command (but not subcommands).
	permissionsRequired []int
	// response is a string that will be sent to the user in response to the command.
	response string
	// embedResponse is a *discordgo.MessageEmbed, if set - has priority over text response.
	embedResponse *discordgo.MessageEmbed
	// description should contain short command description.
	description string
	// usage contains an example of the command usage.
	usage string
}

type iCommand interface {
	Match(sg *Instance, m *discordgo.Message) (matched bool, err error)
	CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error)
	Execute(sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed)

	Trigger() (value string)
	SetTrigger(value string)

	RootOnly() (value bool)
	SetRootOnly(value bool)

	PermissionsRequired() (value []int)
	AddRequiredPermission(value int)

	Response() (value string)
	SetResponse(value string)

	EmbedResponse() (value *discordgo.MessageEmbed)
	SetEmbedResponse(value *discordgo.MessageEmbed)

	Description() (value string)
	SetDescription(value string)

	Usage() (value string)
	SetUsage(value string)
}

func (c *Command) Trigger() (value string) {
	return c.trigger
}

func (c *Command) SetTrigger(value string) {
	c.trigger = value
}

func (c *Command) RootOnly() (value bool) {
	return c.rootOnly
}

func (c *Command) SetRootOnly(value bool) {
	c.rootOnly = value
}

func (c *Command) PermissionsRequired() (value []int) {
	return c.permissionsRequired
}

func (c *Command) AddRequiredPermission(value int) {
	c.permissionsRequired = append(c.permissionsRequired, value)
}

func (c *Command) Response() (value string) {
	return c.response
}

func (c *Command) SetResponse(value string) {
	c.response = value
}

func (c *Command) EmbedResponse() (value *discordgo.MessageEmbed) {
	return c.embedResponse
}

func (c *Command) SetEmbedResponse(value *discordgo.MessageEmbed) {
	c.embedResponse = value
}

func (c *Command) Description() (value string) {
	return c.description
}

func (c *Command) SetDescription(value string) {
	c.description = value
}

func (c *Command) Usage() (value string) {
	return c.usage
}

func (c *Command) SetUsage(value string) {
	c.usage = value
}

func (c *Command) HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed) {
	if c.Trigger() == "" || c.Description() == "" || c.Usage() == "" {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: "Developer of this command did not supply it with description. :frowning:",
			Color:       COLOR_WARNING,
		}
		return
	} else {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: c.Description(),
			Color:       COLOR_INFO,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Usage:",
					Value: fmt.Sprintf("`@%s` %s", sg.Self.Username, c.Usage()),
				},
			},
		}
		return
	}

}

// Match checks if command signature matches the message content.
func (c *Command) Match(sg *Instance, m *discordgo.Message) (matched bool, err error) {
	// By default command is not matched.
	matched = false
	if c.Trigger() == "" && m.Content == "" {
		return true, nil
	}

	if c.Trigger() != "" {
		if strings.HasPrefix(m.Content, c.Trigger()) {
			matched = true
			return true, nil
		}
	}
	return
}

// CheckPermissions checks message author and bot permissions if they match command's required permissions.
func (c *Command) CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error) {
	// By default command is not allowed.
	passed = false

	// For security reasons - every command should have at least one permission set explicitly.
	if len(c.PermissionsRequired()) == 0 {
		err = Error{Text: "Command has no PermissionsRequired[]!"}
		return
	}

	// Get channel to check for permissions.
	channel, err := sg.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Calculate compound permission.
	var compound_perm int = 0
	for _, perm := range c.PermissionsRequired() {
		compound_perm |= perm
	}

	// Make sure bot has the permission required.
	bot_has_perm, err := sg.BotHasPermission(compound_perm, channel)
	if err != nil {
		return
	}
	if !(bot_has_perm) {
		return
	}

	// If user is a root - command is always allowed.
	if sg.IsRoot(m.Author) {
		return true, nil
	}
	// Otherwise if user is not a root a command is root-only - command is not allowed.
	if c.RootOnly() {
		return
	}

	// Make sure user has the permission required.
	user_has_perm, err := sg.UserHasPermission(compound_perm, m.Author, channel)
	if err != nil {
		return
	}
	if !(user_has_perm) {
		return
	}

	// At this time we have checked that:
	// - Command has at least one permission requirement.
	// - Channel we check permissions against exists.
	// - User has all the permissions required.
	// - Bot has all the permissions required.
	// So we can safely say the command is allowed to be executed.
	passed = true
	return
}

// Execute performs commands actions. For basic command it's just a simple text response.
func (c *Command) Execute(sg *Instance, m *discordgo.Message) (err error) {
	if c.embedResponse != nil {
		_, err = c.RespondWithEmbed(sg, m, c.EmbedResponse())
		if err != nil {
			return
		}
	} else {
		_, err = c.RespondWithMention(sg, m, c.Response())
		if err != nil {
			return
		}
	}
	return
}

// Responds to the channel without mention of the original message author.
func (c *Command) Respond(sg *Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

func (c *Command) RespondWithEmbed(sg *Instance, m *discordgo.Message, embed *discordgo.MessageEmbed) (message *discordgo.Message, err error) {
	_, err = sg.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return
	}
	return
}

// Responds to the channel with the original message author mention.
func (c *Command) RespondWithMention(sg *Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	response_text := fmt.Sprintf("%s %s", helpers.UserAsMention(m.Author), text)
	message, err = sg.ChannelMessageSend(m.ChannelID, response_text)
	if err != nil {
		return
	}
	return
}
