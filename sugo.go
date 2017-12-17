// Package sugo is a discord bot framework written in go.
package sugo

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// VERSION contains current version of the Sugo framework.
const VERSION = "0.2.0"

// Instance struct describes bot.
type Instance struct {
	// Bot has everything discordgo.Session has.
	*discordgo.Session
	// root is a user that always has all permissions granted.
	root *discordgo.User
	// Modules contains all Modules loaded by bot.
	Modules []*Module
	// Self contains a giscordgo.User instance of the bot.
	Self *discordgo.User
	// ErrorHandler takes care of errors unhandled elsewhere in the code.
	ErrorHandler func(e error) (err error)
	// DB is literally what it says it is. DataBase.
	DB *sql.DB
	// done is channel that receives Shutdown signals.
	done chan os.Signal
	// triggers contains all the top level triggers for commands.
	triggers []string
}

// Context keys.
type CtxKey string

// Bot contains bot instance.
var Bot = &Instance{}

func init() {
	// Initialize bot Modules list.
	Bot.Modules = []*Module{}
}

// Startup starts the bot up.
func (sg *Instance) Startup(token string, rootUID string) error {
	// Intitialize Shutdown channel.
	sg.done = make(chan os.Signal, 1)

	// Variable to store errors.
	var err error

	// Initialize database.
	sg.DB, err = sql.Open("sqlite3", "./data.sqlite3")
	if err != nil {
		return err
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

	// Get root account info.
	if rootUID != "" {
		root, err := sg.Session.User(rootUID)
		if err != nil {
			return errors.New("Error obtaining root account details... " + err.Error())
		}
		sg.root = root
	}

	// Perform Startup for all Modules.
	for _, module := range sg.Modules {
		if err = module.startup(sg); err != nil {
			return err
		}
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

	return err
}

// Shutdown sends Shutdown signal to the bot's Shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before Shutdown.
func (sg *Instance) teardown() error {
	var err error

	// Perform teardown for all Modules.
	for _, module := range sg.Modules {
		if err = module.teardown(sg); err != nil {
			log.Println(err)
		}
	}

	// Close DB connection.
	sg.DB.Close()

	// Close discord session.
	if err = sg.Session.Close(); err != nil {
		return err
	}
	return nil
}

// triggers is a convenience function to get all top-level commands triggers.
//func (sg *Instance) triggers(m *discordgo.Message) []string {
//	triggers, _ := sg.rootCommand.getSubcommandsTriggers(sg, m)
//	return triggers
//}

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

// GetTriggers returns top level triggers.
func (sg *Instance) GetTriggers() []string {
	return sg.triggers
}

// FindCommand searches for the command in the given Modules, includes all permissions checks.
func (sg *Instance) FindCommand(m *discordgo.Message, q string) (*Command, error) {
	var err error
	var cmd *Command

	// For every module available.
	for _, module := range sg.Modules {
		// Try to find the command in question.
		if cmd, err = module.RootCommand.search(sg, m, q); err != nil {
			return nil, err
		}
		if cmd != nil {
			// Command found.
			return cmd, nil
		}
	}
	// No commands found.
	return nil, nil
}

// onMessageCreate contains all the message processing logic for the bot.
func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	var err error                  // Used to capture and report errors.
	var ctx = context.Background() // Root context.
	var command *Command           // Used to store the command we will execute.
	var q = mc.Content             // Command query string.

	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		Bot.HandleError(errors.New("Bot session error:" + err.Error()))
		Bot.Shutdown()
	}

	// Make sure message author is not a bot.
	if mc.Author.Bot {
		return
	}

	// OnBeforeBotTriggerDetect entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeBotTriggerDetect != nil {
			q, err = module.OnBeforeBotTriggerDetect(Bot, mc.Message, q)
			if err != nil {
				Bot.HandleError(errors.New("OnBeforeMentionDetect error: " + err.Error() + " (" + q + ")"))
			}
		}
	}

	// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
	// for mention detection to work right.
	if strings.HasPrefix(q, "<@!") {
		q = strings.Replace(q, "<@!", "<@", 1)
	}

	// Make sure message starts with bot mention.
	if strings.HasPrefix(strings.TrimSpace(q), Bot.Self.Mention()) {
		// Remove bot trigger from the string.
		q = strings.TrimSpace(strings.TrimPrefix(q, Bot.Self.Mention()))
	} else {
		return
	}

	// Fill context with necessary data.
	// Get Channel.
	channel, err := Bot.ChannelFromMessage(mc.Message)
	if err != nil {
		Bot.HandleError(err)
	}
	// Save into context.
	ctx = context.WithValue(ctx, CtxKey("channel"), channel)

	// Get Guild.
	guild, err := Bot.GuildFromMessage(mc.Message)
	if err != nil {
		Bot.HandleError(err)
	}
	// Save into context.
	ctx = context.WithValue(ctx, CtxKey("guild"), guild)

	// OnBeforeCommandSearch entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeCommandSearch != nil {
			q, err = module.OnBeforeCommandSearch(Bot, mc.Message, q)
			if err != nil {
				Bot.HandleError(errors.New("OnBeforeCommandSearch error: " + err.Error() + " (" + q + ")"))
			}
		}
	}

	// Search for applicable command.
	command, err = Bot.FindCommand(mc.Message, q)
	if err != nil {
		// Unhandled error in command.
		Bot.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + q + ")"))
		Bot.Shutdown()
	}
	if command != nil {
		// Remove command trigger from message string.
		q = strings.TrimSpace(strings.TrimPrefix(q, command.Path()))

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

	Bot.RespondCommandNotFound(mc.Message)

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
		text = "try \"" + c.FullHelpPath(sg) + "\" for details"
	}
	msg, err := sg.RespondDanger(m, "bad command usage\n"+text)
	return msg, err
}

// RespondCommandNotFound responds to the channel with "command not found" message mentioning person that invoked
// command.
func (sg *Instance) RespondCommandNotFound(m *discordgo.Message) (message *discordgo.Message, err error) {
	message, err = sg.RespondDanger(m, "command not found")
	if err != nil {
		return
	}
	return
}

// Respond responds to the channel with an embed without any icons.
func (sg *Instance) Respond(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "@" + m.Author.Username,
		Description: text,
		Color:       ColorDefault,
	})
	return
}

// RespondInfo responds to the channel with the "info" embed.
func (sg *Instance) RespondInfo(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       ":information_source:  @" + m.Author.Username,
		Description: text,
		Color:       ColorInfo,
	})
	return
}

// RespondInfo responds to the channel with the "success" embed.
func (sg *Instance) RespondSuccess(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       ":white_check_mark:  @" + m.Author.Username,
		Description: text,
		Color:       ColorSuccess,
	})
	return
}

// RespondInfo responds to the channel with the "warning" embed.
func (sg *Instance) RespondWarning(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       ":warning: @" + m.Author.Username,
		Description: text,
		Color:       ColorWarning,
	})
	return
}

// RespondInfo responds to the channel with the "Danger" embed.
func (sg *Instance) RespondDanger(m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       ":no_entry: @" + m.Author.Username,
		Description: text,
		Color:       ColorDanger,
	})
	return
}

// HelpEmbed returns help embed for the given command.
func (sg *Instance) HelpEmbed(c *Command, m *discordgo.Message) (embed *discordgo.MessageEmbed, err error) {
	// If command has custom help embed available, return that one.
	if c.HelpEmbed != nil {
		embed, err = c.HelpEmbed(c, sg)
		return
	}
	// Else return automatically generated one.
	embed = c.helpEmbed(sg, m)
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
