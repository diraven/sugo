package sugo

import (
	"strings"
	"github.com/pkg/errors"
)

// Command struct describes basic command type.
type Command struct {
	// Trigger is a sequence of symbols message should start with to match with the command.
	Trigger string
	// Description should contain short command description.
	Description string
	// Usage contains an example of the command usage.
	HasParams bool
	// PermissionsRequired specifies permissions set required by the command.
	PermissionsRequired int
	// Execute code for subcommand.
	Execute func(sg *Instance, r *Request) error
	// SubCommands contains all subcommands of the given command.
	SubCommands []*Command
	// parentCommand contains command, which is parent for this one.
	parent *Command
}

// GetSubcommandsTriggers return all subcommands triggers of the given command available for given user.
func (c *Command) GetSubcommandsTriggers(sg *Instance, r *Request) []string {
	var triggers []string

	// For every subcommand:
	for _, subCommand := range c.SubCommands {
		// If user has permissions to use the command:
		if sg.hasPermissions(r, subCommand.PermissionsRequired) {
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
func (c *Command) match(sg *Instance, r *Request, q string) bool {
	// Ff command is empty and trigger not set - we consider this a match.
	if c.Trigger == "" && q == "" {
		return true
	}

	// If trigger is set and in the query:
	if c.Trigger != "" && strings.HasPrefix(q, c.Trigger) {
		// Make sure user has the permissions necessary to run the command.
		return sg.hasPermissions(r, c.PermissionsRequired)
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
func (c *Command) execute(sg *Instance, req *Request) error {
	if c.Execute != nil {
		err := c.Execute(sg, req)
		return err
	}

	if len(c.SubCommands) > 0 {
		// If there is at least one subcommand and command was not executed - let user know he used command incorrectly.
		_, err := sg.RespondBadCommandUsage(req, "", "")
		return err
	}

	// This should never happen as we validate commands before allowing to add them, but let it be here just in case.
	return errors.New("command has no Execute specified and no subcommands: " + c.GetPath())
}
