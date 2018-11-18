// Package sugo is a discord bot framework written in go.
package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
	"os"
)

// VERSION contains current version of the Instance framework.
const VERSION = "0.6.0"

type RequestMiddleware func(*Request) error
type ResponseMiddleware func(*Response) error
type startupHandler func(sg *Instance) (err error)
type shutdownHandler func(sg *Instance) (err error)

// Instance struct describes bot.
type Instance struct {
	// Trigger specifies what should message start with for the bot to consider it to be command.
	DefaultTrigger string
	// HelpTrigger specifies what should message start with for the bot to consider it to be help command.
	HelpTrigger string
	// Session is a *discordgo.Session bot is wrapped around.
	Session *discordgo.Session
	// Self contains a giscordgo.User instance of the bot.
	Self *discordgo.User
	// RootCommand is a bot root meta-command.
	RootCommand *Command

	// IsTriggered should return true if the bot is to react to command and false otherwise.
	IsTriggered func(req *Request) (triggered bool)
	// ErrorHandler is the function that receives and handles all the errors. Keep in mind that *Request can be nil
	// if error handler is called outside of command request scope.
	ErrorHandler func(req *Request, err error)

	// done is channel that receives Shutdown signals.
	done chan os.Signal
	// startupHandlers are executed sequentially one by one on bot startup.
	startupHandlers []startupHandler
	// shutdownHandlers are executed sequentially one by one on bot shutdown.
	shutdownHandlers []shutdownHandler

	requestMiddlewares  []RequestMiddleware
	responseMiddlewares []ResponseMiddleware
}

// New creates new bot instance.
func New() *Instance {
	// Create our bot.
	sugo := &Instance{}

	// Initialize bot root command.
	sugo.RootCommand = &Command{}

	// Initialize subcommands storage.
	sugo.RootCommand.SubCommands = []*Command{}

	// Return a pointer to our bout.
	return sugo
}

// AddStartupHandler adds function that will be called on bot startup.
func (sg *Instance) AddStartupHandler(handler startupHandler) {
	sg.startupHandlers = append(sg.startupHandlers, handler)
}

// AddShutdownHandler adds function that will be called on bot shutdown.
func (sg *Instance) AddShutdownHandler(handler shutdownHandler) {
	sg.shutdownHandlers = append(sg.shutdownHandlers, handler)
}

// AddRequestMiddleware adds request middleware.
func (sg *Instance) AddRequestMiddleware(m RequestMiddleware) {
	sg.requestMiddlewares = append(sg.requestMiddlewares, m)
}

// AddResponseMiddleware adds response middleware.
func (sg *Instance) AddResponseMiddleware(m ResponseMiddleware) {
	sg.responseMiddlewares = append(sg.responseMiddlewares, m)
}

// AddCommand adds command to the bot's commands list.
func (sg *Instance) AddCommand(c *Command) {
	// Validate command.
	if err := c.validate(); err != nil {
		log.Fatal(err)
	}

	// Set parents for all subcommands.
	c.setParents()

	// Add the command.
	sg.RootCommand.SubCommands = append(sg.RootCommand.SubCommands, c)

	return
}

// hasPermissions calculates if user has all the necessary permissions.
func (sg *Instance) hasPermissions(req *Request, requiredPerms int) (result bool) {
	if requiredPerms != 0 {
		// First of all - get the user perms.
		actualPerms, err := sg.Session.State.UserChannelPermissions(req.Message.Author.ID, req.Channel.ID)
		if err != nil {
			sg.HandleError(nil, errors.Wrap(err, "user permissions retrieval failed"))
			return false
		}

		// Check if user has all the required permissions.
		if (actualPerms | requiredPerms) == actualPerms {
			return true
		}

		// User does not have required permissions.
		return false
	}

	// No permissions specified.
	return true
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
