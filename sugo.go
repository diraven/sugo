// sugo is a discord bot framework written in go.
package sugo

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
)

// VERSION contains current version of the Sugo framework.
const VERSION string = "0.0.21"

// PermissionNone is a permission that is always granted for everybody.
const PermissionNone = 0

// Bot contains bot instance.
var Bot Instance

// Instance interface describes bot.
type Instance struct {
	// Bot has everything discordgo.Session has.
	*discordgo.Session
	// Self contains a giscordgo.User instance of the bot.
	Self *discordgo.User
	// root is a user that always has all permissions granted.
	root *discordgo.User
	// commands contains all the commands loaded into the bot.
	commands map[string]Command
	// data is in-memory data storage.
	data *bot_data
	// CShutdown is channel that receives shutdown signals.
	CShutdown chan os.Signal
}

// Startup starts the bot up.
func (sg *Instance) Startup(token string, root_uid string) (err error) {
	// Intitialize Shutdown channel.
	sg.CShutdown = make(chan os.Signal, 1)

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
	if root_uid != "" {
		root, err := sg.Session.User(root_uid)
		if err != nil {
			// TODO: Report error.
		} else {
			sg.root = root
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

	// Register bot sg.CShutdown channel to receive shutdown signals.
	signal.Notify(sg.CShutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Wait for shutdown signal to arrive.
	<-sg.CShutdown

	// Gracefully shut the bot down.
	sg.teardown()

	return
}

// Shutdown sends shutdown signal to the bot's shutdown channel.
func (sg *Instance) Shutdown() () {
	sg.CShutdown <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before shutdown.
func (sg *Instance) teardown() (err error) {
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

// RegisterCommand adds command to the bot's list of registered commands.
func (sg *Instance) RegisterCommand(trigger string, c Command) (err error) {
	// Initialize commands storage if not yet done.
	if sg.commands == nil {
		sg.commands = make(map[string]Command)
	}

	// Check if given trigger already exists.
	if _, ok := sg.commands[trigger]; ok {
		return Error{
			Text: fmt.Sprintf("Conflicting triggers: Command with top level '%s' trigger already exists.", trigger),
		}
	}

	// Save command into the bot's commands map.
	sg.commands[trigger] = c
	return
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

// UserHasPermission checks if bot has given permission on a given channel.
func (sg *Instance) BotHasPermission(permission int, c *discordgo.Channel) (result bool, err error) {
	result, err = sg.UserHasPermission(permission, sg.Self, c)
	return
}

// onMessageCreate contains all the message processing logic for the bot.
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		// TODO: Report error.
		return
	}

	// Make sure message author is not a bot.
	if m.Author.Bot {
		return
	}

	// Make sure the bot is mentioned in the message, and bot mention is first mention in the message.
	mention := helpers.ConsumeTerm(&m.Content)
	if mention != helpers.UserAsMention(Bot.Self) {
		return
	}

	// Get next term, which is probably command name.
	command_name := helpers.ConsumeTerm(&m.Content)

	if command_name != "" {
		// Command name exists after bot mention.

		// Try to find the command in bot commands.
		command, ok := Bot.commands[command_name]

		if ok {
			// If command has been found.
			is_allowed, err := command.CheckPermissions(&Bot, m.Message)
			if err != nil {
				// TODO: Report error.
			}
			if is_allowed {
				// Execute command.
				err := command.Execute(&Bot, m.Message)
				if err != nil {
					// TODO: Report error.
				}
			}
		} else {
			// Command not found.
			// TODO: Here should probably be something like "Command not found, try help command"
		}
	} else {
		// There is nothing else but bot mention in the message.
		// TODO: Here should be some kind of response where bot presents itself and invites to use "help" command
	}
}
