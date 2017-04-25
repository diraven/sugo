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
	// subCommands contains all subcommands of the given command.
	subCommands []iCommand
	// subCommandsTriggers contains all triggers of subcommands for the help to refer to.
	subCommandsTriggers []string
	// parent_ contains command, which is parent for this one
	parent_ iCommand
}

type iCommand interface {
	Match(sg *Instance, m *discordgo.Message) (matched bool, err error)
	CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error)
	Execute(sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed)
	Path() (value string)

	Trigger() (value string)
	SetTrigger(value string)

	RootOnly() (value bool)
	SetRootOnly(value bool)

	PermissionsRequired() (value []int)
	AddRequiredPermission(value int)

	SubCommands() (value []iCommand)
	AddSubCommand(value iCommand) (err error)

	Response() (value string)
	SetResponse(value string)

	EmbedResponse() (value *discordgo.MessageEmbed)
	SetEmbedResponse(value *discordgo.MessageEmbed)

	Description() (value string)
	SetDescription(value string)

	Usage() (value string)
	SetUsage(value string)
	FullUsage() (value string)

	parent() (command iCommand)
	setParent(command iCommand)
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

func (c *Command) SubCommands() (value []iCommand) {
	return c.subCommands
}

func (c *Command) AddSubCommand(subCommand iCommand) (err error) {
	// Make sure command we are adding was not added anywhere else.
	if subCommand.parent() != nil {
		return Error{fmt.Sprintf("The subcommand is already registered: %s", subCommand)}
	}

	// Set subCommand parent for later reference.
	subCommand.setParent(iCommand(c))

	// Add subCommand.
	c.subCommands = append(c.subCommands, subCommand)

	// Cache subCommand trigger.
	if subCommand.Trigger() != "" {
		c.subCommandsTriggers = append(c.subCommandsTriggers, subCommand.Trigger())
	}
	return nil
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

func (c *Command) Path() (value string) {
	if c.parent_ != nil {
		return fmt.Sprintf("%s %s", c.parent().Path(), c.Trigger())
	} else {
		return fmt.Sprintf("`@%s` %s", Bot.Self.Username, c.Trigger())
	}
}

func (c *Command) FullUsage() (value string) {
	return fmt.Sprintf("%s %s", c.Path(), c.Usage())
}

func (c *Command) parent() (value iCommand) {
	return c.parent_
}

func (c *Command) setParent(value iCommand) {
	c.parent_ = value
}

func (c *Command) HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed) {
	if c.Trigger() == "" || c.Description() == "" {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: "Developer of this command did not supply it with description. :frowning:",
			Color:       COLOR_WARNING,
		}
		return embed
	} else {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: c.Description(),
			Color:       COLOR_INFO,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Usage:",
					Value: c.FullUsage(),
				},
			},
		}
		if len(c.SubCommands()) > 0 {
			embed.Fields = append(embed.Fields,
				&discordgo.MessageEmbedField{
					Name:  "Subcommands:",
					Value: strings.Join(c.subCommandsTriggers, ", "),
				}, &discordgo.MessageEmbedField{
					Name:  "To get help on 'subcommand' type:",
					Value: fmt.Sprintf("`@%s` help %s subcommand", sg.Self.Username, c.Trigger()),
				})
		}
		return embed
	}

}

// Match checks if command signature matches the message content.
func (c *Command) Match(sg *Instance, m *discordgo.Message) (matched bool, err error) {
	// By default command is not matched.
	matched = false

	// If trigger is not set, check if command is empty.
	if c.Trigger() == "" && m.Content == "" {
		return true, nil
	}

	// Trigger is set, see if it's in the message.
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
	if c.EmbedResponse() != nil {
		_, err = c.RespondWithEmbed(sg, m, c.EmbedResponse())
		if err != nil {
			return
		}
	}

	if c.Response() != "" {
		_, err = c.RespondWithMention(sg, m, c.Response())
		if err != nil {
			return
		}
	}

	if len(c.SubCommands()) > 0 {
		_, err = c.RespondWithMention(sg, m, "This command itself does not seem to do anything. Try subcommands instead.")
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
