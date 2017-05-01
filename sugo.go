// Package sugo is a discord bot framework written in go.
package sugo

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// VERSION contains current version of the Sugo framework.
const VERSION string = "0.0.28"

// PermissionNone is a permission that is always granted for everybody.
const PermissionNone = 0

// Instance struct describes bot.
type Instance struct {
	// Bot has everything discordgo.Session has.
	*discordgo.Session
	// Self contains a giscordgo.User instance of the bot.
	Self *discordgo.User
	// root is a user that always has all permissions granted.
	root *discordgo.User
	// rootCommand is the starting point for all the rest of commands.
	rootCommand *Command
	// data is in-memory data storage.
	data *botData
	// done is channel that receives shutdown signals.
	done chan os.Signal
}

// Bot contains bot instance.
var Bot = &Instance{}

func init() {
	// Initialize bot root command, we won't be able to add subcommands to it otherwise.
	Bot.rootCommand = &Command{}
}

// Startup starts the bot up.
func (sg *Instance) Startup(token string, rootUID string) (err error) {
	// Intitialize Shutdown channel.
	sg.done = make(chan os.Signal, 1)

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

	// Perform startup for commands.
	cmdStartup(sg.rootCommand, sg)

	// Register callback for the messageCreate events.
	sg.Session.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = sg.Session.Open()
	if err != nil {
		fmt.Println("Error opening connection... ", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Register bot sg.done channel to receive shutdown signals.
	signal.Notify(sg.done, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal to arrive.
	<-sg.done

	fmt.Println("Termination signal received. Shutting down...")

	// Gracefully shut the bot down.
	sg.teardown()

	fmt.Println("Bye!")

	return
}

// Shutdown sends shutdown signal to the bot's shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before shutdown.
func (sg *Instance) teardown() (err error) {
	// Perform teardown for commands.
	cmdTeardown(sg.rootCommand, sg)

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

// AddCommand is a convenience function to add subcommand to root command.
func (sg *Instance) AddCommand(c *Command) {
	// Save command into the bot's Commands list.
	sg.rootCommand.SubCommands = append(sg.rootCommand.SubCommands, c)
}

// Commands is a convenience function to that returns list of top-level bot commands.
func (sg *Instance) Commands() []*Command {
	return sg.rootCommand.SubCommands
}

// Triggers is a convenience function to get all top-level commands triggers.
func (sg *Instance) Triggers() []string {
	// Save command into the bot's Commands list.
	return sg.rootCommand.subCommandsTriggers
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
func FindCommand(query string, m *discordgo.Message, cmdList []*Command) (output *Command, err error) {
	// For every command in the list provided:
	for _, command := range cmdList {
		// Check if message matches command.
		matched, err := cmdMatch(command, Bot, m)
		if err != nil {
			return nil, err
		}
		if !matched {
			// Message did not match command.
			continue
		}

		// Command matched, check if necessary permissions are present.
		passed, err := cmdCheckPermissions(command, Bot, m)
		if err != nil {
			return nil, err
		}
		if !passed {
			// Message did not pass permissions check.
			return nil, nil
		}

		// Command matched and permissions check passed.

		// Check if there are any subcommands.
		subcommands := command.SubCommands
		if len(subcommands) > 0 {
			// We do have subcommands. Consume original parent command trigger from the message.
			query = strings.TrimSpace(strings.TrimPrefix(query, command.Trigger))

			// Now try to match any of the subcommands.
			subcommand, err := FindCommand(query, m, subcommands)
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
	var err error                                  // Used to capture and report errors.
	var ctx context.Context = context.Background() // Root context.
	var command *Command                           // Used to store the command we will execute.
	var query string = mc.Content                  // Command query string.

	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		log.Fatal(err)
	}

	// Make sure message author is not a bot.
	if mc.Author.Bot {
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
	command, err = FindCommand(query, mc.Message, Bot.rootCommand.SubCommands)
	if err != nil {
		log.Println(err)
	}
	if command != nil {
		// Remove command trigger from message string.
		mc.Content = strings.TrimSpace(strings.TrimPrefix(mc.Content, command.Trigger))

		// And execute command.
		err = cmdExecute(ctx, query, command, Bot, mc.Message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Command not found.
	// TODO: Here should probably be something like a response where bot presents itself and invites to use "help" command.

}

// Execute executes given command.
func (sg *Instance) Execute(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

// RespondText responds to the channel with text without mention of the original message author.
func (sg *Instance) RespondText(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

// RespondEmbed responds to the channel with embed without mention of the original message author.
func (sg *Instance) RespondEmbed(m *discordgo.Message, embed *discordgo.MessageEmbed) (message *discordgo.Message, err error) {
	_, err = sg.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return
	}
	return
}

// RespondTextMention responds to the channel with text with the original message author mention.
func (sg *Instance) RespondTextMention(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	responseText := fmt.Sprintf("%s %s", helpers.UserAsMention(m.Author), text)
	message, err = sg.ChannelMessageSend(m.ChannelID, responseText)
	if err != nil {
		return
	}
	return
}

// HelpEmbed returns automatically generated help embed for the given command.
func (sg *Instance) HelpEmbed(c *Command) (embed *discordgo.MessageEmbed, err error) {
	// If command has custom help embed available, return that one.
	if c.HelpEmbed != nil {
		embed, err = c.HelpEmbed(c, sg)
		return
	}
	// Else return automatically generated one.
	embed = cmdHelpEmbed(c, sg)
	return
}
