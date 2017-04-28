package sugo

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
	"log"
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
	textResponse string
	// embedResponse is a *discordgo.MessageEmbed, if set - has priority over text response.
	embedResponse *discordgo.MessageEmbed
	// description should contain short command description.
	description string
	// usage contains an example of the command usage.
	usage string
	// subCommands contains all subcommands of the given command.
	subCommands []ICommand
	// subCommandsTriggers contains all triggers of subcommands for the help to refer to.
	subCommandsTriggers []string
	// parentCommand contains command, which is parent for this one
	parentCommand ICommand
}

// ICommand describes all the methods every command should have.
type ICommand interface {
	Match(sg *Instance, m *discordgo.Message) (bool, error)
	CheckPermissions(sg *Instance, m *discordgo.Message) (bool, error)
	Execute(sg *Instance, m *discordgo.Message) error
	HelpEmbed(sg *Instance, m *discordgo.Message) *discordgo.MessageEmbed
	Path() string
	Startup() error
	Teardown() error
	startup() error
	teardown() error

	Trigger() string
	SetTrigger(string)

	RootOnly() bool
	SetRootOnly(bool)

	PermissionsRequired() []int
	AddRequiredPermission(int)

	SubCommands() []ICommand
	AddSubCommand(ICommand)
	SubCommandsTriggers() []string

	TextResponse() string
	SetTextResponse(string)

	EmbedResponse() *discordgo.MessageEmbed
	SetEmbedResponse(*discordgo.MessageEmbed)

	Description() string
	SetDescription(string)

	Usage() string
	SetUsage(string)
	FullUsage() string

	parent() ICommand
	setParent(ICommand)
}

// Trigger returns currently set trigger for the command.
func (c *Command) Trigger() (value string) {
	return c.trigger
}

// SetTrigger sets trigger for command.
func (c *Command) SetTrigger(value string) {
	c.trigger = value
}

// RootOnly returns true if command requires user to be root and false otherwise.
func (c *Command) RootOnly() (value bool) {
	return c.rootOnly
}

// SetRootOnly sets root-only requirement for the command.
func (c *Command) SetRootOnly(value bool) {
	c.rootOnly = value
}

// PermissionsRequired returns currently set permissions for the command.
func (c *Command) PermissionsRequired() (value []int) {
	return c.permissionsRequired
}

// AddRequiredPermission adds required permission to the command.
func (c *Command) AddRequiredPermission(value int) {
	c.permissionsRequired = append(c.permissionsRequired, value)
}

// SubCommands returns all registered subcommands of the current command.
func (c *Command) SubCommands() (value []ICommand) {
	return c.subCommands
}

// AddSubCommand adds subcommands to the command.
func (c *Command) AddSubCommand(subCommand ICommand) {
	c.subCommands = append(c.subCommands, subCommand)
}

// SubCommandsTriggers returns triggers for all registered immediate child commands.
func (c *Command) SubCommandsTriggers() []string {
	return c.subCommandsTriggers
}

// TextResponse returns currently set text response of the command.
func (c *Command) TextResponse() (value string) {
	return c.textResponse
}

// SetTextResponse sets text response for the command.
func (c *Command) SetTextResponse(value string) {
	c.textResponse = value
}

// EmbedResponse returns currently set embed response of the command.
func (c *Command) EmbedResponse() (value *discordgo.MessageEmbed) {
	return c.embedResponse
}

// SetEmbedResponse sets embed response for the command.
func (c *Command) SetEmbedResponse(value *discordgo.MessageEmbed) {
	c.embedResponse = value
}

// Description returns currently set description of the command.
func (c *Command) Description() (value string) {
	return c.description
}

// SetDescription sets description for the command.
func (c *Command) SetDescription(value string) {
	c.description = value
}

// Usage returns currently set usage example of the command.
func (c *Command) Usage() (value string) {
	return c.usage
}

// SetUsage sets usage example for the command.
func (c *Command) SetUsage(value string) {
	c.usage = value
}

// startup is internal function called on bot startup.
func (c *Command) startup() error {
	// For every subcommand (if any):
	for _, v := range c.SubCommands() {
		// Build command triggers cache.
		if v.Trigger() != "" {
			c.subCommandsTriggers = append(c.subCommandsTriggers, v.Trigger())
		}

		// Check if command already registered elsewhere.
		if v.parent() != nil {
			return Error{fmt.Sprintf("The subcommand is already registered elsewhere: %s", c.FullUsage())}
		}
		// Set command parent.
		v.setParent(ICommand(c))

		// Run system startup for subcommand.
		v.startup()
	}

	// Run public startup for command.
	c.Startup()

	return nil
}

// teardown is internal function called on bot graceful shutdown.
func (c *Command) teardown() error {
	var err error

	// For every subcommand (if any):
	for _, v := range c.SubCommands() {
		// Run system startup for subcommand.
		err = v.teardown()
		if err != nil {
			log.Printf("Command teardown error: %s\n", err)
		}
	}

	// Run public teardown for command.
	err = c.Teardown()
	if err != nil {
		log.Printf("Command teardown error: %s\n", err)
	}
	return nil
}

// Startup is called on bot startup (may be used for initiating DB connections etc).
func (c *Command) Startup() (err error) {
	return
}

// Teardown is called on graceful bot shutdown (may be used to release resources such as closing files, db connections etc).
func (c *Command) Teardown() (err error) {
	return
}

// Path returns all the triggers from parent commands from outermost to innermost parent.
func (c *Command) Path() (value string) {
	if c.parentCommand != nil {
		return fmt.Sprintf("%s %s", c.parent().Path(), c.Trigger())
	}
	return fmt.Sprintf("`@%s` %s", Bot.Self.Username, c.Trigger())
}

// FullUsage returns full command usage including all parent triggers.
func (c *Command) FullUsage() (value string) {
	return fmt.Sprintf("%s %s", c.Path(), c.Usage())
}

// parent returns parent command if any.
func (c *Command) parent() (value ICommand) {
	return c.parentCommand
}

// setParent sets parent command.
func (c *Command) setParent(value ICommand) {
	c.parentCommand = value
}

// HelpEmbed returns automatically constructed help embed (based on command description, usage etc.) that is ready to
// be sent via discordgo.
func (c *Command) HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed) {
	if c.Trigger() == "" || c.Description() == "" {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: "Developer of this command did not supply it with description. :frowning:",
			Color:       ColorWarning,
		}
		return embed
	}
	embed = &discordgo.MessageEmbed{
		Title:       m.Content,
		Description: c.Description(),
		Color:       ColorInfo,
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
	var compoundPerm int
	for _, perm := range c.PermissionsRequired() {
		compoundPerm |= perm
	}

	// Make sure bot has the permission required.
	botHasPerm, err := sg.BotHasPermission(compoundPerm, channel)
	if err != nil {
		return
	}
	if !(botHasPerm) {
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
	userHasPerm, err := sg.UserHasPermission(compoundPerm, m.Author, channel)
	if err != nil {
		return
	}
	if !(userHasPerm) {
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
		return
	}

	if c.TextResponse() != "" {
		_, err = c.RespondWithMention(sg, m, c.TextResponse())
		if err != nil {
			return
		}
		return
	}

	if len(c.SubCommands()) > 0 {
		_, err = c.RespondWithMention(sg, m, "This command itself does not seem to do anything. Try subcommands instead.")
		if err != nil {
			return
		}
		return
	}

	return
}

// Respond responds to the channel with c.TextResponse text without mention of the original message author.
func (c *Command) Respond(sg *Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

// RespondWithEmbed responds to the chennel with c.Embed embed without mention of the original message author.
func (c *Command) RespondWithEmbed(sg *Instance, m *discordgo.Message, embed *discordgo.MessageEmbed) (message *discordgo.Message, err error) {
	_, err = sg.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return
	}
	return
}

// RespondWithMention responds to the channel with c.TextResponse text with the original message author mention.
func (c *Command) RespondWithMention(sg *Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	responseText := fmt.Sprintf("%s %s", helpers.UserAsMention(m.Author), text)
	message, err = sg.ChannelMessageSend(m.ChannelID, responseText)
	if err != nil {
		return
	}
	return
}
