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
	// parentCommand contains command, which is parent for this one
	parent *Command
}

// GetSubcommandsTriggers return all subcommands triggers available for given user.
func (c *Command) GetSubcommandsTriggers(sg *Instance, r *Request) []string {
	var triggers []string

	// Generate triggers list respecting user permissions.
	for _, subCommand := range c.SubCommands {
		if sg.hasPermissions(r, subCommand.PermissionsRequired) {
			triggers = append(triggers, subCommand.Trigger)
		}
	}
	return triggers
}

// startup is internal function called for each command on bot startup.
func (c *Command) startup(sg *Instance) error {
	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Check if command is already registered elsewhere.
		if v.parent != nil {
			return errors.New("The subcommand is already registered elsewhere: " + c.GetPath())
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

// GetPath returns sequence of triggers from outermost (via the sequence of parents) to the given one.
func (c *Command) GetPath() (value string) {
	if c.parent != nil {
		return strings.TrimSpace(c.parent.GetPath() + " " + c.Trigger)
	}
	return c.Trigger
}

// match is a system matching function that checks if command Trigger matches the start of message content.
func (c *Command) match(sg *Instance, r *Request, q string) bool {
	// If Trigger is not set, check if command is empty.
	if c.Trigger == "" && q == "" {
		return true
	}

	// Trigger is set, see if it's in the message.
	if c.Trigger != "" {
		if strings.HasPrefix(q, c.Trigger) {
			// Now make sure user has the permissions necessary to run the command.
			if !sg.hasPermissions(r, c.PermissionsRequired) {
				return false
			}
			//func (s *State) UserChannelPermissions(userID, channelID string) (apermissions int, err error) {

			return true
		}
	}

	return false
}

// search searches for matching command (including permissions checks) in the given command's subcommands.
func (c *Command) search(sg *Instance, req *Request, q string) (*Command, error) {
	// For every command in subcommands list. We start iterating immediately without considering top level command,
	// because our top level command on bot is an artificial one to contain real ones. So this top level command is
	// simply ignored.
	for _, cmd := range c.SubCommands {
		// Check if message matches command.
		if !cmd.match(sg, req, q) {
			// Message did not match command.
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

	// None of subcommands matched.
	return nil, nil
}

// validate validates commands for them to have either Execute method defined or have subcommands.
func (c *Command) validate() error {
	if c.Execute != nil {
		return nil
	}

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
	for _, subCmd := range c.SubCommands {
		subCmd.parent = c
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

	return errors.New("command has no Execute specified and no subcommands: " + c.GetPath())
}
