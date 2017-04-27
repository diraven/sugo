// Package sugo is a discord bot framework written in go.
package sugo

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// VERSION contains current version of the Sugo framework.
const VERSION string = "0.0.25"

// PermissionNone is a permission that is always granted for everybody.
const PermissionNone = 0

// Bot contains bot instance.
var Bot Instance

// Instance struct describes bot.
type Instance struct {
	// Bot has everything discordgo.Session has.
	*discordgo.Session
	// Self contains a giscordgo.User instance of the bot.
	Self *discordgo.User
	// root is a user that always has all permissions granted.
	root *discordgo.User
	// Commands contains all the Commands loaded into the bot.
	Commands []ICommand
	// Just a cached list of triggers to use with the help command.
	Triggers []string
	// data is in-memory data storage.
	data *botData
	// cShutdown is channel that receives shutdown signals.
	cShutdown chan os.Signal
}

// Startup starts the bot up.
func (sg *Instance) Startup(token string, rootUID string) (err error) {
	// Intitialize Shutdown channel.
	sg.cShutdown = make(chan os.Signal, 1)

	// Initialize data storage.
	_, err = sg.LoadData()
	if err != nil {
		fmt.Println("Error loading data... ", err)
		return
	}

	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session... ", err)
		return
	}

	// Save Discord session into Instance struct.
	sg.Session = s

	// Get bot discordgo.User instance.
	self, err := sg.Session.User("@me")
	if err != nil {
		fmt.Println("Error obtaining account details... ", err)
		return
	}
	sg.Self = self

	// Get root account info.
	if rootUID != "" {
		root, err := sg.Session.User(rootUID)
		if err != nil {
			return err
		}
		sg.root = root
	}

	// Perform startups for all the commands added.
	for _, command := range sg.Commands {
		err = command.Startup()
		if err != nil {
			fmt.Println("Error running startup for command... ", err)
			return
		}
	}

	// Register callback for the messageCreate events.
	sg.Session.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = sg.Session.Open()
	if err != nil {
		fmt.Println("Error opening connection... ", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Register bot sg.cShutdown channel to receive shutdown signals.
	signal.Notify(sg.cShutdown, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal to arrive.
	<-sg.cShutdown

	fmt.Println("Termination signal received. Shutting down...")

	// Gracefully shut the bot down.
	sg.teardown()

	fmt.Println("Bye!")

	return
}

// Shutdown sends shutdown signal to the bot's shutdown channel.
func (sg *Instance) Shutdown() {
	sg.cShutdown <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before shutdown.
func (sg *Instance) teardown() (err error) {
	// Perform teardowns for all the commands added.
	for _, command := range sg.Commands {
		err = command.Teardown()
		if err != nil {
			fmt.Println("Error running teardown for command... ", err)
		}
	}

	// Dump data.
	_, err = sg.DumpData()
	if err != nil {
		return
	}

	// Close discord session.
	err = sg.Session.Close()
	if err != nil {
		return
	}
	return
}

// RegisterCommand adds command to the bot's list of registered Commands.
func (sg *Instance) RegisterCommand(c ICommand) {
	// Save command into the bot's Commands list.
	sg.Commands = append(sg.Commands, c)
	if c.Trigger() != "" {
		sg.Triggers = append(sg.Triggers, c.Trigger())
	}
}

// IsRoot checks if a given user is root.
func (sg *Instance) IsRoot(user *discordgo.User) (result bool) {
	// By default user is not root.
	result = false
	// If root is defined for our bot.
	if sg.root != nil {
		// If root ID is the same as user ID
		if sg.root.ID == user.ID && user.ID != "" {
			// Then the user is root.
			result = true
		}
	}
	return
}

// UserHasPermission checks if given user has given permission on a given channel.
func (sg *Instance) UserHasPermission(permission int, u *discordgo.User, c *discordgo.Channel) (result bool, err error) {
	perms, err := sg.UserChannelPermissions(u.ID, c.ID)
	if err != nil {
		return
	}
	result = (perms | permission) == perms
	return
}

// BotHasPermission checks if bot has given permission on a given channel.
func (sg *Instance) BotHasPermission(permission int, c *discordgo.Channel) (result bool, err error) {
	result, err = sg.UserHasPermission(permission, sg.Self, c)
	return
}

// FindCommand looks for an appropriate (sub)command to execute (taking into account triggers and permissions).
func FindCommand(m *discordgo.Message, cmdList []ICommand) (output ICommand, err error) {
	// For every command in the list provided:
	for _, command := range cmdList {
		// Check if message matches command.
		matched, err := command.Match(&Bot, m)
		if err != nil {
			return nil, err
		}
		if !matched {
			// Message did not match command.
			continue
		}

		// Command matched, check if necessary permissions are present.
		passed, err := command.CheckPermissions(&Bot, m)
		if err != nil {
			return nil, err
		}
		if !passed {
			// Message did not pass permissions check.
			return nil, nil
		}

		// Command matched and permissions check passed.

		// Check if there are any subcommands.
		subcommands := command.SubCommands()
		if len(subcommands) > 0 {
			// We do have subcommands. Consume original parent command trigger from the message.
			m.Content = strings.TrimSpace(strings.TrimPrefix(m.Content, command.Trigger()))

			// Now try to match any of the subcommands.
			subcommand, err := FindCommand(m, subcommands)
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

// onMessageCreate contains all the message processing logic for the bot.
func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		// TODO: Report error.
		return
	}

	// Make sure message author is not a bot.
	if mc.Author.Bot {
		// TODO: Report error.
		return
	}

	// Make sure the bot is mentioned in the message, and bot mention is first mention in the message.
	botMention := helpers.UserAsMention(Bot.Self)
	if strings.HasPrefix(strings.TrimSpace(mc.Content), botMention) {
		// Remove bot mention from the string.
		mc.Content = strings.TrimSpace(strings.TrimPrefix(mc.Content, botMention))
	} else {
		// Bot was not mentioned.
		return
	}

	// Search for applicable command.
	command, err := FindCommand(mc.Message, Bot.Commands)
	if err != nil {
		// TODO: Report error.
	}
	if command != nil {
		err := command.Execute(&Bot, mc.Message)
		if err != nil {
			// TODO: Report error.
		}
	} else {
		// Command not found.
		// TODO: Here should probably be something like a response where bot presents itself and invites to use "help" command.
	}
}
