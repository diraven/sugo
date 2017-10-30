// Package sugo is a discord bot framework written in go.
package sugo

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// VERSION contains current version of the Sugo framework.
const VERSION string = "0.1.3"

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
	// permissionStorage contains struct to get and set per-role command permissions.
	permissions iPermissionsStorage
	// shortcuts contains all the commands shortcuts
	shortcuts iShortcutsStorage
	// Trigger contains global bot trigger (by default it's bot own mention)
	Trigger string
	// done is channel that receives Shutdown signals.
	done chan os.Signal
	// ErrorHandler takes care of unhandled errors.
	ErrorHandler func(e error) (err error)
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

	// Set default permissions storage if one is not specified.
	if sg.permissions == nil {
		sg.permissions = &permissionStorage{}
	}
	if err = sg.permissions.startup(sg); err != nil {
		return
	}

	// Set default shortcuts storage if one is not specified.
	if sg.shortcuts == nil {
		sg.shortcuts = &sShortcutsStorage{}
	}
	if err = sg.shortcuts.startup(sg); err != nil {
		return
	}

	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return errors.New("Error creating Discord session... " + err.Error())
	}

	// Save Discord session into Instance struct.
	sg.Session = s

	// Get bot discordgo.User instance.
	self, err := sg.Session.User("@me")
	if err != nil {
		return errors.New("Error obtaining bot account details... " + err.Error())
	}
	sg.Self = self

	// Set default bot trigger if it's not set beforehand.
	if sg.Trigger == "" {
		// Default trigger is bot's own mention.
		sg.Trigger = sg.Self.Mention() + " "
	}

	// Get root account info.
	if rootUID != "" {
		root, err := sg.Session.User(rootUID)
		if err != nil {
			return errors.New("Error obtaining root account details... " + err.Error())
		}
		sg.root = root
	}

	// Perform Startup for commands.
	if err = sg.rootCommand.startup(sg); err != nil {
		return
	}

	// Register callback for the messageCreate events.
	sg.Session.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	if err = sg.Session.Open(); err != nil {
		return errors.New("Error opening connection... " + err.Error())
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")

	// Register bot sg.done channel to receive Shutdown signals.
	signal.Notify(sg.done, syscall.SIGINT, syscall.SIGTERM)

	// Wait for Shutdown signal to arrive.
	<-sg.done

	// Gracefully shut the bot down.
	err = sg.teardown()

	return
}

// Shutdown sends Shutdown signal to the bot's Shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before Shutdown.
func (sg *Instance) teardown() (err error) {
	// Shutdown permissions storage.
	sg.permissions.teardown(sg)

	// Shutdown shortcuts storage.
	sg.shortcuts.teardown(sg)

	// Perform teardown for commands.
	sg.rootCommand.teardown(sg)

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
//func (sg *Instance) userHasPermission(permission int, c *discordgo.Channel, u *discordgo.User) (result bool, err error) {
//	perms, err := sg.UserChannelPermissions(u.ID, c.ID)
//	if err != nil {
//		return
//	}
//	result = (perms | permission) == perms
//	return
//}

// botHasPermission checks if bot has given permission on a given channel.
//func (sg *Instance) botHasPermission(permission int, c *discordgo.Channel) (result bool, err error) {
//	result, err = sg.userHasPermission(permission, c, sg.Self)
//	return
//}

// onMessageCreate contains all the message processing logic for the bot.
func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	var err error                                  // Used to capture and report errors.
	var ctx context.Context = context.Background() // Root context.
	var command *Command                           // Used to store the command we will execute.
	var q string = mc.Content                      // Command query string.

	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		Bot.HandleError(errors.New("Bot session error:" + err.Error()))
		Bot.Shutdown()
	}

	// Make sure message author is not a bot.
	if mc.Author.Bot {
		return
	}

	// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
	// for mention detection to work right.
	if strings.HasPrefix(q, "<@!") {
		q = strings.Replace(q, "<@!", "<@", 1)
	}

	// Make sure message starts with bot trigger.
	if strings.HasPrefix(strings.TrimSpace(q), Bot.Trigger) {
		// Remove bot trigger from the string.
		q = strings.TrimSpace(strings.TrimPrefix(q, Bot.Trigger))
	} else {
		return
	}

	// Process shortcuts.
	for _, shortcut := range Bot.shortcuts.all() {
		if strings.Index(q, shortcut.Short) == 0 {
			if len(q) == len(shortcut.Short) || string(q[len(shortcut.Short)]) == " " {
				q = strings.Replace(q, shortcut.Short, shortcut.Long, 1)
				break
			}
		}
	}

	// Search for applicable command.
	command, err = Bot.rootCommand.search(Bot, q, mc.Message)
	if err != nil {
		// Unhandled error in command.
		Bot.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + q + ")"))
		Bot.Shutdown()
	}
	if command != nil {
		// Remove command trigger from message string.
		q = strings.TrimSpace(strings.TrimPrefix(q, command.path()))

		// And execute command.
		err = command.execute(ctx, q, Bot, mc.Message)
		if err != nil {
			if strings.Contains(err.Error(), "\"code\": 50013") {
				// Insufficient permissions, bot configuration issue.
				Bot.HandleError(errors.New("Bot permissions error: " + err.Error() + " (" + q + ")"))
			} else {
				// Other discord errors.
				Bot.HandleError(errors.New("Bot command execute error: " + err.Error() + " (" + q + ")"))
				Bot.Shutdown()
			}
		}
		return
	}

	//Bot.respondCommandNotFound(mc.Message)

	// Command not found.
}

// RespondText responds to the channel with text without mention of the original message author.
func (sg *Instance) RespondText(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

// RespondBadCommandUsage responds to the channel with "incorrect command usage" message mentioning person that invoked
// command.
// _, err = sg.RespondBadCommandUsage(m, c)
func (sg *Instance) RespondBadCommandUsage(m *discordgo.Message, c *Command, text string) (message *discordgo.Message, err error) {
	if text == "" {
		text = "Try \"" + c.FullHelpPath(sg) + "\" for details."
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Incorrect command usage.",
		Description: text,
		Color:       ColorDanger,
	}

	message, err = sg.RespondEmbed(m, embed)
	if err != nil {
		return
	}
	return
}

// respondCommandNotFound responds to the channel with "command not found" message mentioning person that invoked
// command.
func (sg *Instance) respondCommandNotFound(m *discordgo.Message) (message *discordgo.Message, err error) {
	embed := &discordgo.MessageEmbed{
		Title: "Command not found.",
		Color: ColorDanger,
	}

	message, err = sg.RespondEmbed(m, embed)
	if err != nil {
		return
	}
	return
}

// RespondEmbed responds to the channel with embed without mention of the original message author.
func (sg *Instance) RespondEmbed(m *discordgo.Message, embed *discordgo.MessageEmbed) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}

// RespondTextMention responds to the channel with text with the original message author mention.
func (sg *Instance) RespondTextMention(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	responseText := m.Author.Mention() + " " + text
	message, err = sg.ChannelMessageSend(m.ChannelID, responseText)
	return
}

// RespondSuccessMention responds to the channel with white checkmark on a green background with the original message author mention.
func (sg *Instance) RespondSuccessMention(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.RespondTextMention(m, ":white_check_mark: "+text)
	return
}

// RespondFailMention responds to the channel with white checkmark on a green background with the original message author mention.
func (sg *Instance) RespondFailMention(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	if text == "" {
		text = "Oops... Something went wrong!"
	}
	message, err = sg.RespondTextMention(m, ":exclamation: "+text)
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
	embed = c.helpEmbed(sg)
	return
}

// ChannelFromMessage returns a *discordgo.Channel struct from given *discordgo.Message struct.
func (sg *Instance) ChannelFromMessage(m *discordgo.Message) (c *discordgo.Channel, err error) {
	return sg.State.Channel(m.ChannelID)
}

// GuildFromMessage returns a *discordgo.Guild struct from given *discordgo.Message struct.
func (sg *Instance) GuildFromMessage(m *discordgo.Message) (g *discordgo.Guild, err error) {
	c, err := sg.ChannelFromMessage(m)
	if err != nil {
		return
	}
	return sg.State.Guild(c.GuildID)
}

// MemberFromMessage returns a *discordgo.Member struct from given *discordgo.Message struct.
func (sg *Instance) MemberFromMessage(m *discordgo.Message) (mr *discordgo.Member, err error) {
	g, err := sg.GuildFromMessage(m)
	if err != nil {
		return
	}
	return sg.State.Member(g.ID, m.Author.ID)
}

// HandleError handles unexpected errors that were returned unhandled elsewhere.
func (sg *Instance) HandleError(e error) (err error) {
	if sg.ErrorHandler != nil {
		err = sg.ErrorHandler(e)
	} else {
		log.Println(e)
		sg.Shutdown()
	}
	return
}
