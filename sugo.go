// Package sugo is a discord bot framework written in go.
package sugo

import (
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3" // We do not need actual sqlite driver reference in package.
	"os"
	"github.com/pkg/errors"
	"log"
)

// VERSION contains current version of the Instance framework.
const VERSION = "0.4.0"

// Instance struct describes bot.
type Instance struct {
	// Trigger specifies what should message start with for the bot to consider it to be command.
	Trigger string
	// HelpTrigger specifies what should message start with for the bot to consider it to be help command.
	HelpTrigger string
	// Session is a *discordgo.Session bot is wrapped around.
	Session *discordgo.Session
	// Self contains a giscordgo.User instance of the bot.
	Self *discordgo.User
	// RootCommand is a bot root meta-command.
	RootCommand *Command
	// root is a user whose ID was specified during bot error.
	root *discordgo.User
	// done is channel that receives Shutdown signals.
	done chan os.Signal
	// aliases contains all bot command aliases.
	aliases *aliases

	// errorHandlers are executed sequentially one by one on bot error.
	startupHandlers []func(sg *Instance) error
	// shutdownHandlers are executed sequentially one by one on bot shutdown.
	shutdownHandlers []func(sg *Instance) error
	// shutdownHandlers are executed sequentially one by one on bot shutdown.
	errorHandlers []func(err error)
}

// New creates new bot instance.
func New() *Instance {
	// Create our bot.
	sugo := &Instance{}

	// Initialize bot root command.
	sugo.RootCommand = &Command{}

	// Initialize aliases storage.
	sugo.aliases = &aliases{}

	// Initialize subcommands storage.
	sugo.RootCommand.SubCommands = []*Command{}

	// Return a pointer to our bout.
	return sugo
}

// AddStartupHandler adds function that will be called on bot startup.
func (sg *Instance) AddStartupHandler(handler func(sg *Instance) error) {
	sg.startupHandlers = append(sg.startupHandlers, handler)
}

// AddShutdownHandler adds function that will be called on bot shutdown.
func (sg *Instance) AddShutdownHandler(handler func(sg *Instance) error) {
	sg.shutdownHandlers = append(sg.shutdownHandlers, handler)
}

// AddErrorHandler adds function that will be called on error if it's unhandled elsewhere.
func (sg *Instance) AddErrorHandler(handler func(err error)) {
	sg.errorHandlers = append(sg.errorHandlers, handler)
}

// AddCommand adds command to the bot's commands list.
func (sg *Instance) AddCommand(c *Command) {
	// Validate command.
	if err := c.validate(); err != nil {
		sg.HandleError(err)
	}

	// Set parents for all subcommands.
	c.setParents()

	// Add the command.
	sg.RootCommand.SubCommands = append(sg.RootCommand.SubCommands, c)
}

// hasPermissions calculates if user has all the necessary permissions.
func (sg *Instance) hasPermissions(r *Request, requiredPerms int) bool {
	// First of all - get the user perms.
	actualPerms, err := sg.Session.State.UserChannelPermissions(r.Message.Author.ID, r.Channel.ID)
	if err != nil {
		sg.HandleError(errors.Wrap(err, "user actual permissions retrieval failed"))
		return false
	}

	// Check if user has all the required permissions.
	if (actualPerms | requiredPerms) == actualPerms {
		return true
	}

	// User does not have required permissions.
	return false
}

// FindCommand searches for the command in the modules registered.
func (sg *Instance) FindCommand(req *Request, q string) (*Command, error) {
	var err error
	var cmd *Command

	// Try to find the command in question.
	if cmd, err = sg.RootCommand.search(sg, req, q); err != nil {
		return nil, err
	}
	if cmd != nil {
		// Command found.
		return cmd, nil
	}

	// No commands found.
	return nil, nil
}

// HandleError handles unexpected errors that were returned unhandled elsewhere.
func (sg *Instance) HandleError(e error) {
	// If there are any error handlers registered:
	if len(sg.errorHandlers) > 0 {
		// Run error handlers.
		for _, handler := range sg.errorHandlers {
			handler(e)
		}
	}

	// Otherwise just put error into the log.
	log.Printf("%+v", e)
}
