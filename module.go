package sugo

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Module struct {
	RootCommand *Command
	Startup     func(sg *Instance) error
	Teardown    func(sg *Instance) error

	// OnPermissionsCheck is a function that is called every time permissions checked for a command,
	// possible return values:
	// false - command is considered denied
	// true - command is considered allowed
	// nil - use default command permissions
	OnPermissionsCheck func(sg *Instance, c *Command, m *discordgo.Message) (*bool, error)

	// OnBeforeCommandSearch is called before command search is performed, but after query string is prepared
	// returned value replaces query string that will be used for command search.
	OnBeforeCommandSearch func(sg *Instance, m *discordgo.Message, q string) (string, error)

	// OnBeforeCommandSearch is called before query is tested for bot mention.
	OnBeforeBotTriggerDetect func(sg *Instance, m *discordgo.Message, q string) (string, error)

	// onPresenceUpdate happens every time member presence is updated for guild.
	OnPresenceUpdate func(sg *Instance, pu *discordgo.PresenceUpdate) error
}

// startup is internal function called for each module on bot startup.
func (m *Module) startup(sg *Instance) error {
	// For Modules with commands - fill parent fields and triggers cache as well as validate triggers.
	if m.RootCommand != nil {
		// if trigger is set:
		if m.RootCommand.Trigger != "" {
			// Make sure trigger is unique.
			for _, v := range sg.triggers {
				if v == m.RootCommand.Trigger {
					return errors.New("trigger is already registered: " + m.RootCommand.Trigger)
				}
			}

			// Add top level trigger to the bot triggers cache.
			sg.triggers = append(sg.triggers, m.RootCommand.Trigger)
		}

		// For every subcommand (if any):
		for _, v := range m.RootCommand.SubCommands {
			// Check if command is already registered elsewhere.
			if v.parent != nil {
				return errors.New("The subcommand is already registered elsewhere: " + m.RootCommand.Path())
			}
			// Set command parent.
			v.parent = m.RootCommand

			// Run system startup for subcommand.
			if err := v.startup(sg); err != nil {
				return err
			}
		}
	}

	// Run custom startup for module if set.
	if m.Startup != nil {
		if err := m.Startup(sg); err != nil {
			return err
		}
	}

	return nil
}

// teardown is internal function called for each module on bot graceful Shutdown.
func (m *Module) teardown(sg *Instance) error {
	// Run public teardown for module if set.
	if m.Teardown != nil {
		if err := m.Teardown(sg); err != nil {
			return err
		}
	}
	return nil
}
