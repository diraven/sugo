// Package sugo is a discord bot framework written in go.
package sugo

import (
	"database/sql"
	"errors"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3" // We do not need actual sqlite driver reference in package.
	"github.com/nicksnyder/go-i18n/i18n"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	ErrorHandler func(error) error
	// DB is literally what it says it is. DataBase.
	DB *sql.DB
	// done is channel that receives Shutdown signals.
	done chan os.Signal
	// triggers contains all the top level triggers for commands.
	triggers []string
}

// Bot contains bot instance.
var Bot = &Instance{}

func init() {
	// Initialize bot Modules list.
	Bot.Modules = []*Module{}

	// Make translations dir.
	const translationsPath = "translations"
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		err = os.Mkdir(translationsPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Fill translation file with default data.
	if err := ioutil.WriteFile(translationsPath+"/en.yaml", []byte(`- id: program_greeting
  translation: Hello world`), 0655); err != nil {
		log.Fatal(err)
	}

	// Load all translation files.
	files, err := filepath.Glob(translationsPath + "/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		i18n.MustLoadTranslationFile(v)
	}
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
func (sg *Instance) FindCommand(req *Request, q string) (*Command, error) {
	var err error
	var cmd *Command

	// For every module available.
	for _, module := range sg.Modules {
		// Try to find the command in question.
		if cmd, err = module.RootCommand.search(sg, req, q); err != nil {
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
	}

	log.Println(e)
	sg.Shutdown()
	return nil
}

// processMessage processes given message.
func (sg *Instance) processMessage(m *discordgo.Message) {
	var err error // Used to capture and report errors.
	var request = &Request{}

	// Save message.
	request.Message = m

	// Save query.
	request.Query = m.Content

	// Get and save Channel.
	request.Channel, err = sg.State.Channel(m.ChannelID)

	// Get and save Guild if available.
	if request.Channel.GuildID != "" {
		request.Guild, err = sg.State.Guild(request.Channel.GuildID)
		if err != nil {
			Bot.HandleError(err)
		}
	}

	// Get translation function.
	tfunc, err := i18n.Tfunc("en")
	if err != nil {
		Bot.HandleError(err)
	}
	request.TranslateFunc = &tfunc

	// OnMessageCreate entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnMessageCreate != nil {
			if err = module.OnMessageCreate(Bot, request); err != nil {
				Bot.HandleError(errors.New("OnMessageCreate error: " + err.Error()))
			}
		}
	}

	// Ignore any message that is coming from bot.
	if m.Author.Bot {
		return
	}

	if request.Channel.Type == discordgo.ChannelTypeDM {
		// It's Direct Messaging channel. Every message here is in fact a direct message to the bot, so we consider
		// it to be command without further checks.

	}

	if request.Channel.Type == discordgo.ChannelTypeGuildText || request.Channel.Type == discordgo.ChannelTypeGroupDM {
		// It's either Guild Text channel or multiple people direct group channel.
		// In order to detect command we need to account for trigger.

		// OnBeforeBotTriggerDetect entry point for Modules.
		for _, module := range Bot.Modules {
			if module.OnBeforeBotTriggerDetect != nil {
				module.OnBeforeBotTriggerDetect(Bot, request)
				if err != nil {
					Bot.HandleError(errors.New("OnBeforeMentionDetect error: " + err.Error() + " (" + request.Query + ")"))
				}
			}
		}

		// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
		// for mention detection to work right.
		if strings.HasPrefix(request.Query, "<@!") {
			request.Query = strings.Replace(request.Query, "<@!", "<@", 1)
		}

		// Make sure message starts with bot mention.
		if strings.HasPrefix(strings.TrimSpace(request.Query), Bot.Self.Mention()) {
			// Remove bot trigger from the string.
			request.Query = strings.TrimSpace(strings.TrimPrefix(request.Query, Bot.Self.Mention()))
		} else {
			return
		}

	}

	// OnBeforeCommandSearch entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeCommandSearch != nil {
			err = module.OnBeforeCommandSearch(Bot, request)
			if err != nil {
				Bot.HandleError(errors.New("OnBeforeCommandSearch error: " + err.Error() + " (" + request.Query + ")"))
			}
		}
	}

	// Search for applicable command.
	request.Command, err = Bot.FindCommand(request, request.Query)
	if err != nil {
		// Unhandled error in command.
		Bot.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + request.Query + ")"))
	}

	if request.Command != nil {
		// Remove command trigger from message string.
		request.Query = strings.TrimSpace(strings.TrimPrefix(request.Query, request.Command.Path()))

		// Make sure command is possible to execute (i.e. it supports DM if channel is of DM type).
		if !request.Command.AllowDM && request.Channel.Type != discordgo.ChannelTypeGuildText {
			return
		}

		// And execute command.
		err = request.Command.execute(Bot, request)
		if err != nil {
			if strings.Contains(err.Error(), "\"code\": 50013") {
				// Insufficient permissions, bot configuration issue.
				Bot.HandleError(errors.New("Bot permissions error: " + err.Error() + " (" + request.Query + ")"))
			} else {
				// Other discord errors.
				Bot.HandleError(errors.New("Bot command execute error: " + err.Error() + " (" + request.Query + ")"))
			}
			Bot.HandleError(err)
		}
	}

	// Command not found.
}
