// Package sugo is a discord bot framework written in go.
package sugo

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// VERSION contains current version of the Sugo framework.
const VERSION string = "0.1.0"

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
	// Debug determines if bot is in the debug mode (false by default).
	Debug bool
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

	// Set default permissions storage if one is not specified.
	sg.DebugLog(0, "Initiating permissions...")
	if sg.permissions == nil {
		sg.permissions = &permissionStorage{}
	}
	sg.permissions.startup(sg)
	sg.DebugLog(0, "Done.")

	// Set default shortcuts storage if one is not specified.
	sg.DebugLog(0, "Initiating shortcuts...")
	if sg.shortcuts == nil {
		sg.shortcuts = &shortcutsStorage{}
	}
	sg.shortcuts.startup(sg)
	sg.DebugLog(0, "Done.")

	// Create a new Discord session using the provided bot token.
	sg.DebugLog(0, "Initiating Discord session...")
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("sError creating Discord session... ", err)
		return
	}
	sg.DebugLog(0, "Done.")

	// Save Discord session into Instance struct.
	sg.Session = s

	// Get bot discordgo.User instance.
	sg.DebugLog(0, "Getting bot user...")
	self, err := sg.Session.User("@me")
	if err != nil {
		log.Println("sError obtaining account details... ", err)
		return
	}
	sg.Self = self
	sg.DebugLog(0, "Done. ID:", sg.Self.ID)

	sg.DebugLog(0, "Getting root info...")
	// Get root account info.
	if rootUID != "" {
		root, err := sg.Session.User(rootUID)
		if err != nil {
			return err
		}
		sg.root = root
	}
	sg.DebugLog(0, "Done. Root:", sg.root.Username)

	// Perform Startup for commands.
	sg.DebugLog(0, "Performing commands startup...")
	sg.rootCommand.startup(sg)
	sg.DebugLog(0, "Done.")

	// Register callback for the messageCreate events.
	sg.DebugLog(0, "Registering onMessageCreate callback...")
	sg.Session.AddHandler(onMessageCreate)
	sg.DebugLog(0, "Done.")

	// Open the websocket and begin listening.
	sg.DebugLog(0, "Opening socket...")
	err = sg.Session.Open()
	if err != nil {
		log.Println("sError opening connection... ", err)
		return
	}
	sg.DebugLog(0, "Done.")
	log.Println("Bot is now running. Press CTRL-C to exit.")

	// Register bot sg.done channel to receive Shutdown signals.
	signal.Notify(sg.done, syscall.SIGINT, syscall.SIGTERM)

	// Wait for Shutdown signal to arrive.
	<-sg.done

	sg.DebugLog(0, "Termination signal received. Shutting down...")
	// Gracefully shut the bot down.
	sg.teardown()
	sg.DebugLog(0, "Done.")

	return
}

// Shutdown sends Shutdown signal to the bot's Shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before Shutdown.
func (sg *Instance) teardown() (err error) {
	// Shutdown permissions storage.
	sg.DebugLog(1, "Tearing down permissions...")
	sg.permissions.teardown(sg)
	sg.DebugLog(1, "Done.")

	// Shutdown permissions storage.
	sg.DebugLog(1, "Tearing down shortcuts...")
	sg.shortcuts.teardown(sg)
	sg.DebugLog(1, "Done.")

	// Perform teardown for commands.
	sg.DebugLog(1, "Tearing down commands...")
	sg.rootCommand.teardown(sg)
	sg.DebugLog(1, "Done.")

	// Close discord session.
	sg.DebugLog(1, "Closing Discord session...")
	err = sg.Session.Close()
	sg.DebugLog(1, "Done.")
	if err != nil {
		return
	}
	return
}

// AddCommand is a convenience function to add subcommand to root command.
func (sg *Instance) AddCommand(c *Command) {
	sg.DebugLog(0, "Adding command:", c.path())
	// Save command into the bot's commands list.
	sg.rootCommand.SubCommands = append(sg.rootCommand.SubCommands, c)
	sg.DebugLog(0, "Done.")
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
		Bot.Shutdown()
		log.Fatalln("ERROR:", err)
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
	if Bot.Trigger != "" && strings.HasPrefix(strings.TrimSpace(q), Bot.Trigger) {
		// Remove bot trigger from the string.
		q = strings.TrimSpace(strings.TrimPrefix(q, Bot.Trigger))
	} else {
		// There is either no trigger or its not in the message. Check if we have our bot mentioned.
		if strings.HasPrefix(strings.TrimSpace(q), Bot.Self.Mention()) {
			// Remove bot mention from the string.
			q = strings.TrimSpace(strings.TrimPrefix(q, Bot.Self.Mention()))
		} else {
			// Bot was not mentioned.
			return
		}
	}

	// Process shortcuts.
	Bot.DebugLog(0, "Looking for shortcuts...")
	for shortcut, command := range Bot.shortcuts.all() {
		log.Println(q)
		if strings.Index(q, shortcut) == 0 {
			if len(q) == len(shortcut) || string(q[len(shortcut)]) == " " {
				q = strings.Replace(q, shortcut, command, 1)
				Bot.DebugLog(0, "Shortcut found! Replaced \"", shortcut, "\" with \"", command, "\"")
				break
			}
		}
	}
	Bot.DebugLog(0, "Done.")

	// Search for applicable command.
	command, err = Bot.rootCommand.search(Bot, q, mc.Message)
	if err != nil {
		Bot.Shutdown()
		log.Fatalln("ERROR:", err)
	}
	if command != nil {
		Bot.DebugLog(0, "Got command to execute:", command.path())
		// Remove command trigger from message string.
		q = strings.TrimSpace(strings.TrimPrefix(q, command.path()))

		// And execute command.
		Bot.DebugLog(0, "Executing...")
		err = command.execute(ctx, q, Bot, mc.Message)
		if err != nil {
			if strings.Contains(err.Error(), "\"code\": 50013") {
				// Insufficient permissions, bot configuration issue.
				log.Println("ERROR:", err)
			} else {
				Bot.Shutdown()
				log.Fatalln("ERROR:", err)
			}
		}
		Bot.DebugLog(0, "Command execution finished:", command.path())
		return
	}

	Bot.respondCommandNotFound(mc.Message)

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
func (sg *Instance) RespondBadCommandUsage(m *discordgo.Message, c *Command) (message *discordgo.Message, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       "Incorrect command usage.",
		Description: "Try \"" + c.FullHelpPath(sg) + "\" for details",
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

// DebugLog puts vars into the log if bot debug is enabled.
func (sg *Instance) DebugLog(nesting int, v ...interface{}) {
	if nesting > 0 {
		prefix := make([]interface{}, nesting, nesting)
		for i := 0; i < nesting; i++ {
			prefix[i] = ">"
		}
		v = append(prefix, v...)
	}
	if sg.Debug {
		log.Println(v...)
	}
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
