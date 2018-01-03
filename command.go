package sugo

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Command struct describes basic command type.
type Command struct {
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
	// HasParams specifies if command is allowed to have additional parameters after the command path itself.
	ParamsAllowed bool
	// DMAble specifies if command can be used via direct messages.
	DMAble bool

	// parentCommand contains command, which is parent for this one
	parent *Command

	// Custom execute code for subcommand.
	Execute func(sg *Instance, c *Command, m *discordgo.Message, q string) error
}

// getSubcommandsTriggers return all subcommands triggers available for given user.
func (c *Command) GetSubcommandsTriggers(sg *Instance, m *discordgo.Message) ([]string, error) {
	var triggers []string

	// Generate triggers list respecting user permissions.
	for _, subCommand := range c.SubCommands {
		command, err := sg.FindCommand(m, subCommand.Path())
		if err != nil {
			return triggers, err
		}
		if command != nil {
			triggers = append(triggers, subCommand.Trigger)
		}
	}
	return triggers, nil
}

// startup is internal function called for each command on bot startup.
func (c *Command) startup(sg *Instance) error {
	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Check if command is already registered elsewhere.
		if v.parent != nil {
			return errors.New("The subcommand is already registered elsewhere: " + c.Path())
		}
		// Set command parent.
		v.parent = c

		// Run system startup for subcommand.
		if err := v.startup(sg); err != nil {
			return err
		}
	}

	return nil
}

// teardown is internal function called for each command on bot graceful Shutdown.
func (c *Command) teardown(sg *Instance) error {
	// !!!! Here be some internal code to tear commands down... some day. May be.

	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Run system teardown for subcommand.
		if err := v.teardown(sg); err != nil {
			return err
		}
	}

	return nil
}

// Path returns sequence of triggers from outermost to innermost command for the given one.
func (c *Command) Path() (value string) {
	if c.parent != nil {
		return strings.TrimSpace(c.parent.Path() + " " + c.Trigger)
	}
	return c.Trigger
}

func (c *Command) FullHelpPath(sg *Instance) (value string) {
	return "help " + c.Path()
}

// match is a system matching function that checks if command trigger matches the start of message content.
func (c *Command) match(sg *Instance, m *discordgo.Message, q string) bool {
	// If trigger is not set, check if command is empty.
	if c.Trigger == "" && q == "" {
		return true
	}

	// Trigger is set, see if it's in the message.
	if c.Trigger != "" {
		if strings.HasPrefix(q, c.Trigger) {
			return true
		}
	}
	return false
}

// search searches for matching command (including permissions checks) in the given command's subcommands.
func (c *Command) search(sg *Instance, m *discordgo.Message, q string) (*Command, error) {
	// Check if message matches command.
	if !c.match(sg, m, q) {
		// Message did not match command.
		return nil, nil
	}

	// Command matched, check if necessary permissions are present.
	passed, err := c.checkPermissions(sg, m)
	if err != nil {
		return nil, err
	}
	if !passed {
		// Message did not pass permissions check.
		return nil, nil
	}

	// Command matched and permissions check passed!
	// Consume original parent command trigger from the message.
	q = strings.TrimSpace(strings.TrimPrefix(q, c.Trigger))

	// Check if there are any subcommands.
	if len(c.SubCommands) > 0 {
		// We do have subcommands.
		for _, subCommand := range c.SubCommands {
			// Now try to match any of the subcommands.
			result, err := subCommand.search(sg, m, q)
			if err != nil {
				return nil, err
			}
			// If we were able to get subcommand that matched:
			if result != nil {
				return result, nil
			}
		}
	}

	// Either there are no subcommands, or none of those worked. Return parent command, but only if it has no params
	// or params allowed.
	if c.ParamsAllowed || q == "" {
		return c, nil
	}

	// We did not find a command that satisfies all the requirements.
	return nil, nil
}

// checkPermissions checks if given user has necessary permissions to use the command. The function is called
// sequentially for topmost command and following the path to the subcommand in question.
func (c *Command) checkPermissions(sg *Instance, m *discordgo.Message) (bool, error) {
	// If user is a root - command is always allowed.
	if sg.isRoot(m.Author) {
		return true, nil
	}

	// Otherwise if user is not a root and command is root-only - command is not allowed.
	if c.RootOnly {
		return false, nil
	}

	// Get channel.
	channel, err := sg.State.Channel(m.ChannelID)
	if err != nil {
		return false, err
	}
	// Check if we should ignore the command because it's disabled for default channel.
	if !c.AllowDefaultChannel && channel.ID == channel.GuildID {
		return false, nil
	}

	// Now check if we have any additional Modules handling permission checks and use those.
	var allowedFound bool // specifies if any of the Modules returned explicit result

	for _, module := range sg.Modules {
		if module.OnPermissionsCheck != nil {
			passed, err := module.OnPermissionsCheck(sg, c, m)
			if err != nil {
				// In case of error - return error and deny command.
				return false, err
			}
			if passed == nil { // Undefined return.
				continue // Just go on to the next module.
			}
			if *passed { // Command explicitly allowed
				// Mark the fact we have found module that allows the command and go on to the next ones.
				allowedFound = true
				continue
			} else { // Command explicitly disallowed.
				return false, nil // We return false if at least one module disallowed the command execution.
			}
		}
	}

	if allowedFound { // At this point if there are Modules found with explicit return - they did allow the command.
		return true, nil
	}

	// There are no special permissions set for any of the Modules. Just back to default.
	return c.PermittedByDefault, nil
}

// execute is a default command execution function.
func (c *Command) execute(q string, sg *Instance, m *discordgo.Message) error {
	var actionPerformed bool

	if c.Execute != nil {
		// Run custom command Execute if set.

		if err := c.Execute(sg, c, m, q); err != nil {
			return err
		}
		actionPerformed = true
	}

	if c.TextResponse != "" {
		// Send command text response if set.

		if _, err := sg.Respond(m, "", c.TextResponse, ColorPrimary, ""); err != nil {
			return err
		}
		actionPerformed = true
	}

	if c.EmbedResponse != nil {
		// Send command embed response if set.

		if _, err := sg.ChannelMessageSendEmbed(m.ChannelID, c.EmbedResponse); err != nil {
			return err
		}
		actionPerformed = true
	}

	if !actionPerformed {
		if len(c.SubCommands) > 0 {
			// If there is at least one subcommand and no other actions taken - explain it to the user.
			_, err := sg.RespondBadCommandUsage(
				m, c, "", "")
			return err
		}

		// We did nothing and there are no subcommands...
		_, err := Bot.RespondInfo(m, "", "looks like this command just does nothing... what is it here for anyways?")
		return err
	}

	return nil
}
