package sugo

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
	// IgnoreDefaultChannel specifies if command works in guild default channel.
	AllowDefaultChannel bool
	// PermittedByDefault specifies if command is allowed to be used by default. Default is false.
	PermittedByDefault bool
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

	Execute   func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed func(c *Command, sg *Instance) (embed *discordgo.MessageEmbed, err error)
	Startup   func(c *Command, sg *Instance) (err error)
	Teardown  func(c *Command, sg *Instance) (err error)
}

func (c *Command) getSubcommandsTriggers(sg *Instance, m *discordgo.Message) (triggers []string, err error) {
	triggers = []string{}

	// Generate triggers list respecting user permissions.
	for _, trigger := range c.subCommandsTriggers {
		command, err := c.search(sg, m, trigger)
		if err != nil {
			return triggers, err
		}
		if command != nil {
			triggers = append(triggers, trigger)
		}
	}
	return triggers, nil
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
			return errors.New("The subcommand is already registered elsewhere: " + c.path())
		}
		// Set command parent.
		v.parent = c

		// Run system startup for subcommand.
		if err = v.startup(sg); err != nil {
			return
		}
	}

	// Run custom startup for command if set.
	if c.Startup != nil {
		if err = c.Startup(c, sg); err != nil {
			return
		}
	}

	return
}

// teardown is internal function called for each command on bot graceful Shutdown.
func (c *Command) teardown(sg *Instance) (err error) {
	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Here be some internal code to tear commands down... some day. May be.

		// Run system teardown for subcommand.
		if err = v.teardown(sg); err != nil {
			return
		}
	}

	// Run public teardown for command if set.
	if c.Teardown != nil {
		if err = c.Teardown(c, sg); err != nil {
			return
		}
	}
	return
}

// path returns sequence of triggers from outermost to innermost command for the given one.
func (c *Command) path() (value string) {
	if c.parent != nil {
		return strings.TrimSpace(c.parent.path() + " " + c.Trigger)
	}
	return c.Trigger
}

func (c *Command) FullHelpPath(sg *Instance) (value string) {
	return sg.Trigger + "help " + c.path()
}

// fullUsage returns full command usage including all parent triggers.
func (c *Command) fullUsage(sg *Instance) (value string) {
	return sg.Trigger + c.path() + " " + c.Usage
}

// helpEmbed is a default implementation of help embed builder.
func (c *Command) helpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed) {
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
	// Get subcommands triggers respecting user permissions.
	subcommandsTriggers, _ := c.getSubcommandsTriggers(sg, m)


	if len(c.SubCommands) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Subcommands:",
				Value: strings.Join(subcommandsTriggers, ", "),
			}, &discordgo.MessageEmbedField{
				Name:  "To get help on 'subcommand' type:",
				Value: fmt.Sprintf("`@%s` help %s subcommand", sg.Self.Username, c.Trigger),
			})
	}
	return embed

}

// match is a system matching function that checks if command trigger matches the start of message content.
func (c *Command) match(sg *Instance, m *discordgo.Message, q string) (matched bool, err error) {
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

// search searches for matching command (including permissions checks) in the given command's subcommands.
func (c *Command) search(sg *Instance, m *discordgo.Message, q string) (output *Command, err error) {
	// For every command in the list provided:
	for _, command := range c.SubCommands {

		// Check if message matches command.
		matched, err := command.match(sg, m, q)
		if err != nil {
			return nil, err
		}
		if !matched {
			// Message did not match command.
			continue
		}

		// Command matched, check if necessary permissions are present.
		passed, err := command.checkPermissions(sg, m)
		if err != nil {
			return nil, err
		}
		if !passed {
			// Message did not pass permissions check.
			return nil, nil
		}

		// Command matched and permissions check passed.

		// Check if there are any subcommands.
		if len(command.SubCommands) > 0 {
			// We do have subcommands. Consume original parent command trigger from the message.
			q = strings.TrimSpace(strings.TrimPrefix(q, command.Trigger))

			// Now try to match any of the subcommands.
			subcommand, err := command.search(sg, m, q)
			if err != nil {
				return nil, err
			}
			// If we were able to get subcommand that matched, return it.
			if subcommand != nil {
				return subcommand, nil
			}
		}

		// Either there are no subcommands, or none of those worked. Return parent command.
		return command, nil
	}
	// No commands matched.
	return nil, nil
}

// checkCheckPermissions checks if given user has necessary permissions to use the command. The function is called
// sequentially for topmost command and following the path to the subcommand in question.
func (c *Command) checkPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error) {
	// If user is a root - command is always allowed.
	if sg.isRoot(m.Author) {
		return true, nil
	}

	// Otherwise if user is not a root and command is root-only - command is not allowed.
	if c.RootOnly {
		return
	}

	// Get channel.
	channel, err := sg.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	// Check if we should ignore the command because it's disabled for default channel.
	if !c.AllowDefaultChannel && channel.ID == channel.GuildID {
		return // passed=false, err=nil
	}

	// Get guild member.
	member, err := sg.State.Member(channel.GuildID, m.Author.ID)
	if err != nil {
		return
	}

	// Now we need to check if we have any settings for every role user has sequentially starting from the topmost one.

	// For each role user has.
	var role *discordgo.Role
	var position int = 0 // position of the custom role setting found
	var found bool       // if custom role setting found

	// Start with checking "everyone" role permissions.
	role, err = sg.State.Role(channel.GuildID, channel.GuildID)
	if err != nil {
		return false, err // Just make sure we are safe and return false in case of any errors.
	}
	isAllowed, exists := sg.permissions.get(sg, c.path(), role.ID)
	if exists {
		found = true
		passed = isAllowed
	}

	// And now check all the rest of the user roles.
	for _, roleID := range member.Roles {
		// Get role itself.
		role, err = sg.State.Role(channel.GuildID, roleID)
		if err != nil {
			return false, err // Just make sure we are safe and return false in case of any errors.
		}
		// Check if role is allowed to use the command.
		isAllowed, exists := sg.permissions.get(sg, c.path(), role.ID)

		// If custom role setting exists and it's position less then the one we have already found (role that is higher
		// takes precedence over the lower ones):
		if exists && role.Position >= position {
			position = role.Position // Store position of the role.
			found = true             // Make sure we know we have found a custom setting.
			passed = isAllowed       // Update return value.
		}
	}

	if found {
		// If we have found the custom role setting - we just return what we have found.
		return
	}

	// There are no special permissions set for any of the user's roles. Fall back to default.
	passed = c.PermittedByDefault
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
