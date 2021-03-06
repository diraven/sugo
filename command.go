package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"strings"
)

// Command struct describes basic command type.
type Command struct {
	// Trigger is a sequence of symbols message should start with to match with the command.
	Trigger string
	// Description should contain short command description.
	Description string
	// HasParams specifies if command can have additional parameters in Request string.
	HasParams bool
	// PermissionsRequired specifies permissions set required by the command.
	PermissionsRequired int
	// RequireGuild specifies if this command works in guild chats only.
	RequireGuild bool
	// Execute method is executed if Request string matches the given command.
	Execute func(req *Request) (*Response, error)
	// SubCommands contains all subcommands of the given command.
	SubCommands []*Command
	// parentCommand contains command, which is parent for this one.
	parent *Command
}

// GetSubcommandsTriggers return all subcommands triggers of the given command available for given user.
func (c *Command) GetSubcommandsTriggers(sg *Instance, req *Request) []string {
	var triggers []string

	// For every subcommand:
	for _, subCommand := range c.SubCommands {
		// If user has permissions to use the command:
		if sg.hasPermissions(req, subCommand.PermissionsRequired) {
			// Add subcommand trigger to the list.
			triggers = append(triggers, subCommand.Trigger)
		}
	}

	return triggers
}

// GetPath returns sequence of triggers from outermost (via the sequence of parents) to the given one.
func (c *Command) GetPath() (value string) {
	if c.parent != nil {
		return strings.TrimSpace(c.parent.GetPath() + " " + c.Trigger)
	}
	return c.Trigger
}

// match is a system matching function that checks if command Trigger matches the start of message content.
func (c *Command) match(sg *Instance, req *Request, q string) bool {
	// If command is for guild Text channels only and executed elsewhere - it's not a match.
	if c.RequireGuild && req.Channel.Type != discordgo.ChannelTypeGuildText {
		return false
	}

	// If command is empty and trigger not set - we consider this a match.
	if c.Trigger == "" && q == "" {
		return true
	}

	// If trigger is set and in the query:
	if c.Trigger != "" && strings.HasPrefix(q, c.Trigger) {
		// Make sure user has permissions necessary to run the command.
		return sg.hasPermissions(req, c.PermissionsRequired)
	}

	// If no trigger is set and query is not empty then it's not a match.
	return false
}

// search searches for matching command (including permissions checks) in the given command's subcommands.
func (c *Command) search(sg *Instance, req *Request, q string) (*Command, error) {
	// For every command in subcommands list. We start iterating immediately without considering top level command,
	// because our top level command on bot is an artificial one to contain real ones. So this top level command is
	// simply ignored.
	for _, cmd := range c.SubCommands {
		// If message does not match command:
		if !cmd.match(sg, req, q) {
			// Continue searching.
			continue
		}

		// Make sure to strip away the Trigger of the parent command we have already found as matching.
		q = strings.TrimSpace(strings.TrimPrefix(q, cmd.Trigger))

		// Try to find subcommand that matches the remainder of the query.
		subCmd, err := cmd.search(sg, req, q)
		if err != nil {
			return nil, err
		}
		if subCmd != nil {
			// If we found matching subcommand, return it.
			return subCmd, nil
		}

		// Otherwise return our parent command whose subcommands we were iterating over.
		// Either q should be empty (fully consumed by matching) or the command we are going to return should be able to
		// accept and process parameters.
		// It's done to exclude false positives that tend to happen when you try to use subcommands and spell them
		// improperly, which results in a situation where we return parent command with it's improperly spelled
		// subcommand Trigger as a parameter.
		if q == "" || cmd.HasParams {
			return cmd, nil
		}

		// Otherwise continue with searching another command.
	}

	// No subcommands matched.
	return nil, nil
}

// validate validates commands for them to have either Execute method defined or have subcommands.
func (c *Command) validate() error {
	// If command has Execute function defined - we consider it valid and subcommands do not matter.
	if c.Execute != nil {
		return nil
	}

	// If there is no Execute function, but there is at least one subcommand:
	if len(c.SubCommands) > 0 {
		// Perform validation recursively for every subcommand.
		for _, subCmd := range c.SubCommands {
			return subCmd.validate()
		}
	}

	return errors.New("command has neither subcommands nor Execute method defined: " + c.GetPath())
}

// setParents sets parents for all commands for easier reference.
func (c *Command) setParents() {
	// For every subcommand:
	for _, subCmd := range c.SubCommands {
		// Set subcommand parent to current command.
		subCmd.parent = c
		// Do the same recursively for all sub-subcommands.
		subCmd.setParents()
	}
}

// execute is a default command execution function.
func (c *Command) execute(sg *Instance, req *Request) (resp *Response, err error) {
	// There is always either execute defined or subcommands available as enforced by validate() on command add.

	// If execute method defined - use it.
	if c.Execute != nil {
		return c.Execute(req)
	}

	// Otherwise there must be subcommands. Notify user that command is used incorrectly.
	resp = req.NewResponse(ResponseDanger, "", "I'm unable to execute this command itself, try subcommands instead")

	return
}
