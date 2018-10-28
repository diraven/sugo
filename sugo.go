// Package sugo is a discord bot framework written in go.
package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

// VERSION contains current version of the Instance framework.
const VERSION = "0.5.3"

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
	// done is channel that receives Shutdown signals.
	done chan os.Signal
	// aliases contains all bot command aliases.
	aliases *aliases

	// IsTriggered should return true if the bot is to react to command and false otherwise.
	IsTriggered func(req *Request) (triggered bool, err error)
	// errorHandlers are executed sequentially one by one on bot error.
	startupHandlers []func(sg *Instance) (err error)
	// shutdownHandlers are executed sequentially one by one on bot shutdown.
	shutdownHandlers []func(sg *Instance) (err error)
	// shutdownHandlers are executed sequentially one by one on bot shutdown.
	errorHandlers []func(
		err error,
		req *Request,
	)
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
func (sg *Instance) AddErrorHandler(handler func(
	err error,
	req *Request,
)) {
	sg.errorHandlers = append(sg.errorHandlers, handler)
}

// AddCommand adds command to the bot's commands list.
func (sg *Instance) AddCommand(c *Command) (err error) {
	// Validate command.
	if err = c.validate(); err != nil {
		return
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
			sg.HandleError(errors.Wrap(err, "user permissions retrieval failed"), req)
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

// HandleError handles unexpected errors that were returned unhandled elsewhere.
func (sg *Instance) HandleError(
	err error,
	req *Request,
) {
	// If there are any error handlers registered:
	if len(sg.errorHandlers) > 0 {
		// Run error handlers.
		for _, handler := range sg.errorHandlers {
			handler(err, req)
		}
		return
	}

	// Otherwise just put error into the log.
	log.Println(err)
}

func (sg *Instance) isTriggered(req *Request) (triggered bool, err error) {
	// Use custom IsTriggered function is provided.
	if sg.IsTriggered != nil {
		return sg.IsTriggered(req)
	}

	if req.Channel.Type == discordgo.ChannelTypeDM {
		// It's Direct Messaging Channel. Every message here is in fact a direct message to the bot, so we consider
		// it to be command without any further checks for prefixes.
		return true, nil
	} else if req.Channel.Type == discordgo.ChannelTypeGuildText || req.Channel.Type == discordgo.ChannelTypeGroupDM {
		// It's either Guild Text Channel or multiple people direct group Channel.
		// In order to detect command we need to check for bot Triggereq.

		// If bot Trigger is set and command starts with that Trigger:
		if sg.DefaultTrigger != "" && strings.HasPrefix(req.Query, sg.DefaultTrigger) {
			// Replace custom Trigger with bot mention for it to be detected as bot Triggereq.
			req.Query = strings.Replace(req.Query, sg.DefaultTrigger, req.Sugo.Self.Mention(), 1)
		}

		// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
		// for mention detection to work right.
		if strings.HasPrefix(req.Query, "<@!") {
			req.Query = strings.Replace(req.Query, "<@!", "<@", 1)
		}

		// If the message starts with bot mention:
		if strings.HasPrefix(strings.TrimSpace(req.Query), req.Sugo.Self.Mention()) {
			// Remove bot Trigger from the string.
			req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Sugo.Self.Mention()))
			// Bot is triggered.
			return true, nil
		}

		// Otherwise bot is not triggered.
		return false, nil
	}

	// We ignore all other channel types and consider bot not triggered.
	return false, nil
}
