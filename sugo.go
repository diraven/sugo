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
const VERSION string = "0.0.29"

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
	// done is channel that receives Shutdown signals.
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

	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("sError creating Discord session... ", err)
		return
	}

	// Save Discord session into Instance struct.
	sg.Session = s

	// Get bot discordgo.User instance.
	self, err := sg.Session.User("@me")
	if err != nil {
		fmt.Println("sError obtaining account details... ", err)
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

	// Perform Startup for commands.
	cmdStartup(sg.rootCommand, sg)

	// Register callback for the messageCreate events.
	sg.Session.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = sg.Session.Open()
	if err != nil {
		fmt.Println("sError opening connection... ", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Register bot sg.done channel to receive Shutdown signals.
	signal.Notify(sg.done, syscall.SIGINT, syscall.SIGTERM)

	// Wait for Shutdown signal to arrive.
	<-sg.done

	fmt.Println("Termination signal received. Shutting down...")

	// Gracefully shut the bot down.
	sg.teardown()

	fmt.Println("Bye!")

	return
}

// Shutdown sends Shutdown signal to the bot's Shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before Shutdown.
func (sg *Instance) teardown() (err error) {
	// Perform teardown for commands.
	cmdTeardown(sg.rootCommand, sg)

	// Close discord session.
	err = sg.Session.Close()
	if err != nil {
		return
	}
	return
}

// AddCommand is a convenience function to add subcommand to root command.
func (sg *Instance) AddCommand(c *Command) {
	// Save command into the bot's commands list.
	sg.rootCommand.SubCommands = append(sg.rootCommand.SubCommands, c)
}

// commands is a convenience function to that returns list of top-level bot commands.
func (sg *Instance) commands() []*Command {
	return sg.rootCommand.SubCommands
}

// triggers is a convenience function to get all top-level commands triggers.
func (sg *Instance) triggers() []string {
	return sg.rootCommand.subCommandsTriggers
}

// isRoot checks if a given user is root.
func (sg *Instance) isRoot(user *discordgo.User) (result bool) {
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

// userHasPermission checks if given user has given permission on a given channel.
func (sg *Instance) userHasPermission(permission int, c *discordgo.Channel, u *discordgo.User) (result bool, err error) {
	perms, err := sg.UserChannelPermissions(u.ID, c.ID)
	if err != nil {
		return
	}
	result = (perms | permission) == perms
	return
}

// botHasPermission checks if bot has given permission on a given channel.
func (sg *Instance) botHasPermission(permission int, c *discordgo.Channel) (result bool, err error) {
	result, err = sg.userHasPermission(permission, c, sg.Self)
	return
}

// findCommand looks for an appropriate (sub)command to execute (taking into account triggers and permissions).
func findCommand(q string, m *discordgo.Message, cmdList []*Command) (output *Command, err error) {
	// For every command in the list provided:
	for _, command := range cmdList {
		// Check if message matches command.
		matched, err := cmdMatch(command, q, Bot, m)
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
			q = strings.TrimSpace(strings.TrimPrefix(q, command.Trigger))

			// Now try to match any of the subcommands.
			subcommand, err := findCommand(q, m, subcommands)
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
	var q string = mc.Content                      // Command query string.

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
	if strings.HasPrefix(strings.TrimSpace(q), botMention) {
		// Remove bot mention from the string.
		q = strings.TrimSpace(strings.TrimPrefix(q, botMention))
	} else {
		// Bot was not mentioned.
		return
	}

	// Search for applicable command.
	command, err = findCommand(q, mc.Message, Bot.rootCommand.SubCommands)
	if err != nil {
		log.Println(err)
	}
	if command != nil {
		// Remove command trigger from message string.
		q = strings.TrimSpace(strings.TrimPrefix(q, cmdPath(command)))

		// And execute command.
		err = cmdExecute(ctx, q, command, Bot, mc.Message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Command not found.
	// TODO: Here should probably be something like a response where bot presents itself and invites to use "help" command.

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

// helpEmbed returns automatically generated help embed for the given command.
func (sg *Instance) helpEmbed(c *Command) (embed *discordgo.MessageEmbed, err error) {
	// If command has custom help embed available, return that one.
	if c.HelpEmbed != nil {
		embed, err = c.HelpEmbed(c, sg)
		return
	}
	// Else return automatically generated one.
	embed = cmdHelpEmbed(c, sg)
	return
}
