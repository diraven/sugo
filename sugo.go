// Package sugo is a discord bot framework written in go.
package sugo

import (
	"log"
	"os"
	"github.com/bwmarrin/discordgo"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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
