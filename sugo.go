// Package sugo is a discord bot framework written in go.
package sugo

import (
	"log"
	"os"
	"github.com/bwmarrin/discordgo"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"context"
	"errors"
)

// VERSION contains current version of the Sugo framework.
const VERSION = "0.2.1"

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
	ErrorHandler func(error) error
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

// HandleError handles unexpected errors that were returned unhandled elsewhere.
func (sg *Instance) HandleError(e error) error {
	if sg.ErrorHandler != nil {
		return sg.ErrorHandler(e)
	} else {
		log.Println(e)
		sg.Shutdown()
	}
	return nil
}

// processMessage processes given message.
func (sg *Instance) processMessage(m *discordgo.Message) error {
	var err error                  // Used to capture and report errors.
	var ctx = context.Background() // Root context.
	var command *Command           // Used to store the command we will execute.
	var q = m.Content              // Command query string.

	// OnMessageCreate entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnMessageCreate != nil {
			if err = module.OnMessageCreate(Bot, m); err != nil {
				Bot.HandleError(errors.New("OnMessageCreate error: " + err.Error()))
			}
		}
	}

	// Make sure message author is not a bot.
	if m.Author.Bot {
		return nil
	}

	// OnBeforeBotTriggerDetect entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeBotTriggerDetect != nil {
			q, err = module.OnBeforeBotTriggerDetect(Bot, m, q)
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
		return nil
	}

	// Fill context with necessary data.
	// Get Channel.
	channel, err := Bot.ChannelFromMessage(m)
	if err != nil {
		Bot.HandleError(err)
	}
	// Save into context.
	ctx = context.WithValue(ctx, CtxKey("channel"), channel)

	// Get Guild.
	guild, err := Bot.GuildFromMessage(m)
	if err != nil {
		Bot.HandleError(err)
	}
	// Save into context.
	ctx = context.WithValue(ctx, CtxKey("guild"), guild)

	// OnBeforeCommandSearch entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeCommandSearch != nil {
			q, err = module.OnBeforeCommandSearch(Bot, m, q)
			if err != nil {
				Bot.HandleError(errors.New("OnBeforeCommandSearch error: " + err.Error() + " (" + q + ")"))
			}
		}
	}

	// Search for applicable command.
	command, err = Bot.FindCommand(m, q)
	if err != nil {
		// Unhandled error in command.
		Bot.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + q + ")"))
		Bot.Shutdown()
	}
	if command != nil {
		// Remove command trigger from message string.
		q = strings.TrimSpace(strings.TrimPrefix(q, command.Path()))

		// And execute command.
		err = command.execute(ctx, q, Bot, m)
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
		return nil
	}

	// Command not found.
	return nil
}
