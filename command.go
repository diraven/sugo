package sugo

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
	"log"
	"strings"
	"time"
)

// Command struct describes basic command type.
type Command struct {
	// Timeout
	Timeout time.Duration
	// Trigger is a sequence of symbols message should start with to match with the command.
	Trigger string
	// RootOnly determines if the command is supposed to be used by root only.
	RootOnly bool
	// PermissionsRequired is a slice of all permissions required by the command (but not subcommands).
	Permissions []int
	// Response is a string that will be sent to the user in response to the command.
	TextResponse string
	// EmbedResponse is a *discordgo.MessageEmbed, if set - has priority over text response.
	EmbedResponse *discordgo.MessageEmbed
	// Description should contain short command description.
	Description string
	// Usage contains an example of the command usage.
	Usage string
	// SubCommands contains all subcommands of the given command.
	SubCommands []*Command

	// parentCommand contains command, which is parent for this one
	parent *Command
	// subCommandsTriggers contains all triggers of subcommands for the help to refer to.
	subCommandsTriggers []string

	Execute   func(ctx context.Context, c *Command, query string, sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed func(c *Command, sg *Instance) (embed *discordgo.MessageEmbed, err error)
	Startup   func(c *Command, sg *Instance) (err error)
	Teardown  func(c *Command, sg *Instance) (err error)
}

// startup is internal function called for each command on bot startup.
func (c *Command) startup(sg *Instance) (err error) {
	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Build command triggers cache.
		if v.Trigger != "" {
			c.subCommandsTriggers = append(c.subCommandsTriggers, v.Trigger)
		}

		// Check if command is already registered elsewhere.
		if v.parent != nil {
			return sError{
				fmt.Sprintf("The subcommand is already registered elsewhere: %s", c.path()),
			}
		}
		// Set command parent.
		v.parent = c

		// Run system startup for subcommand.
		v.startup(sg)
	}

	// Run public startup for command if set.
	if c.Startup != nil {
		err = c.Startup(c, sg)
		if err != nil {
			return
		}
	}

	return
}

// teardown is internal function called for each command on bot graceful Shutdown.
func (c *Command) teardown(sg *Instance) error {
	var err error

	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Here be some internal code to tear commands down... some day. May be.

		// Run system startup for subcommand.
		err = v.teardown(sg)
		if err != nil {
			log.Fatal("Command teardown error: ", err)
		}
	}

	// Run public teardown for command if set.
	if c.Teardown != nil {
		err = c.Teardown(c, sg)
		if err != nil {
			return sError{fmt.Sprintf("Command custom teardown error: %s\n", err)}
		}
	}
	return nil
}

// path returns all the triggers from parent commands from outermost to innermost parent.
func (c *Command) path() (value string) {
	if c.parent != nil {
		return strings.TrimSpace(c.parent.path() + " " + c.Trigger)
	}
	return c.Trigger
}

func (c *Command) FullHelpPath(sg *Instance) (value string) {
	return helpers.UserAsMention(sg.Self) + " help " + c.path()
}

// fullUsage returns full command usage including all parent triggers.
func (c *Command) fullUsage(sg *Instance) (value string) {
	return helpers.UserAsMention(sg.Self) + " " + c.path() + " " + c.Usage
}

// helpEmbed is a default implementation of help embed builder.
func (c *Command) helpEmbed(sg *Instance) (embed *discordgo.MessageEmbed) {
	embed = &discordgo.MessageEmbed{
		Title:       c.path(),
		Description: c.Description,
		Color:       ColorInfo,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Usage:",
				Value: c.fullUsage(sg),
			},
		},
	}
	if len(c.SubCommands) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Subcommands:",
				Value: strings.Join(c.subCommandsTriggers, ", "),
			}, &discordgo.MessageEmbedField{
				Name:  "To get help on 'subcommand' type:",
				Value: fmt.Sprintf("`@%s` help %s subcommand", sg.Self.Username, c.Trigger),
			})
	}
	return embed

}

// match is a system matching function that checks if command trigger matches the start of message content.
func (c *Command) match(q string, sg *Instance, m *discordgo.Message) (matched bool, err error) {
	// By default command is not matched.
	matched = false

	// If trigger is not set, check if command is empty.
	if c.Trigger == "" && q == "" {
		return true, nil
	}

	// Trigger is set, see if it's in the message.
	if c.Trigger != "" {
		if strings.HasPrefix(q, c.Trigger) {
			matched = true
			return true, nil
		}
	}
	return
}

// checkCheckPermissions checks message author and bot permissions if they match command's required permissions.
func (c *Command) checkPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error) {
	// By default command is not allowed.
	passed = false

	// For security reasons - every command should have at least one permission set explicitly.
	if len(c.Permissions) == 0 {
		err = sError{"Command has no Permissions[]!"}
		return
	}

	// Get channel to check for permissions.
	channel, err := sg.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Calculate compound permission.
	var compoundPerm int
	for _, perm := range c.Permissions {
		compoundPerm |= perm
	}

	// Make sure bot has the permission required.
	botHasPerm, err := sg.botHasPermission(compoundPerm, channel)
	if err != nil {
		return
	}
	if !(botHasPerm) {
		return
	}

	// If user is a root - command is always allowed.
	if sg.isRoot(m.Author) {
		return true, nil
	}
	// Otherwise if user is not a root a command is root-only - command is not allowed.
	if c.RootOnly {
		return
	}

	// Make sure user has the permission required.
	userHasPerm, err := sg.userHasPermission(compoundPerm, channel, m.Author)
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

// execute is a default command execution function.
func (c *Command) execute(ctx context.Context, q string, sg *Instance, m *discordgo.Message) (err error) {
	var actionPerformed bool

	// Set timeout to the context if requested by user.
	if c.Timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	if c.Execute != nil {
		// Run custom command Execute if set.
		err = c.Execute(ctx, c, q, sg, m)
		if err != nil {
			return
		}
		actionPerformed = true
	}

	if c.TextResponse != "" {
		// Send command text response if set.
		_, err = sg.RespondTextMention(m, c.TextResponse)
		if err != nil {
			return
		}
		actionPerformed = true
	}

	if c.EmbedResponse != nil {
		// Send command embed response if set.
		_, err = sg.RespondEmbed(m, c.EmbedResponse)
		if err != nil {
			return
		}
		actionPerformed = true
	}

	if !actionPerformed {
		if len(c.SubCommands) > 0 {
			// If there is at least one subcommand and no other actions taken - explain it to the user.
			_, err = sg.RespondTextMention(
				m,
				"This command itself does not seem to do anything. Try "+c.FullHelpPath(sg)+".",
			)
			return
		}

		// We did nothing and there are no subcommands...
		_, err = Bot.RespondTextMention(m, "Looks like this command just does nothing... What is it here for?")
		return
	}

	return
}
