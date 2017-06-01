package sugo

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type iShortcutsStorage interface {
	all() (shortcuts map[string]string)
	get(sg *Instance, shortcut string) (command string, exists bool)
	set(sg *Instance, shortcut string, command string)
	del(sg *Instance, shortcut string)
	load(sg *Instance) (data_length int, err error)
	save(sg *Instance) (data_length int, err error)
	startup(sg *Instance) error
	teardown(sg *Instance) error
}

const SHORTCUTS_DATA_FILENAME = "shortcuts.json"

type shortcutsStorage struct {
	Shortcuts map[string]string
}

func (p *shortcutsStorage) all() (shortcuts map[string]string) {
	return p.Shortcuts
}

func (p *shortcutsStorage) get(sg *Instance, shortcut string) (command string, exists bool) {
	command, exists = p.Shortcuts[shortcut]
	return
}

func (p *shortcutsStorage) set(sg *Instance, shortcut string, command string) {
	p.Shortcuts[shortcut] = command
}

func (p *shortcutsStorage) del(sg *Instance, shortcut string) {
	delete(p.Shortcuts, shortcut)
}

func (p *shortcutsStorage) load(sg *Instance) (data_length int, err error) {
	if _, error_type := os.Stat(SHORTCUTS_DATA_FILENAME); os.IsNotExist(error_type) {
		sg.DebugLog(0, "No shortcuts file found. Empty storage initialized.")
		// File to load data from does not exist.
		// Check if perms storage is empty and initialize it.
		shortcutsStorage := sg.shortcuts.(*shortcutsStorage)
		if shortcutsStorage.Shortcuts == nil {
			shortcutsStorage.Shortcuts = make(map[string]string)
		}
		return
	}

	// Load file.
	data, err := ioutil.ReadFile(SHORTCUTS_DATA_FILENAME)
	if err != nil {
		return
	}

	// Decode JSON data.
	json.Unmarshal(data, sg.shortcuts.(*shortcutsStorage))
	if err != nil {
		return
	}

	// Log the operation results.
	data_length = len(data)
	sg.DebugLog(0, "Shortcuts loaded successfully,", data_length, "bytes read.")

	return
}

func (p *shortcutsStorage) save(sg *Instance) (data_length int, err error) {
	// Encode our data into JSON.
	data, err := json.Marshal(sg.shortcuts.(*shortcutsStorage))
	if err != nil {
		return
	}

	// Save data into file.
	err = ioutil.WriteFile(SHORTCUTS_DATA_FILENAME, data, 0644)
	if err != nil {
		return
	}

	data_length = len(data)
	sg.DebugLog(0, "Shortcuts saved successfully,", data_length, "bytes written.")

	return
}

func (p *shortcutsStorage) startup(sg *Instance) (err error) {
	_, err = p.load(sg)
	return
}

func (p *shortcutsStorage) teardown(sg *Instance) (err error) {
	_, err = p.save(sg)
	return
}

// Help shows help section for appropriate command.
var CmdShortcuts = &Command{
	Trigger:     "shortcuts",
	RootOnly:    true,
	Description: "Allows to manipulate shortcuts. Lists all shortcuts.",
	Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
		var result string = ""

		for shortcut, command := range sg.shortcuts.(*shortcutsStorage).Shortcuts {
			result = result + shortcut + " -> " + command + "\n"
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Currently configured shortcuts",
			Description: result,
		}

		sg.RespondEmbed(m, embed)
		return
	},
	SubCommands: []*Command{
		{
			Trigger:     "load",
			Description: "Loads shortcuts from storage.",
			Usage:       "shortcuts save",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				data_length, err := sg.shortcuts.load(sg)
				if err != nil {
					return
				}
				sg.RespondTextMention(m, "Shortcuts loaded. "+strconv.FormatInt(int64(data_length), 10)+" bytes read.")
				return
			},
		},
		{
			Trigger:     "save",
			Description: "Saves shortcuts to storage.",
			Usage:       "shortcuts load",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				data_length, err := sg.shortcuts.save(sg)
				if err != nil {
					return
				}
				sg.RespondTextMention(m, "Shortcuts saved. "+strconv.FormatInt(int64(data_length), 10)+" bytes written.")
				return
			},
		},
		{
			Trigger:     "set",
			Description: "Adds new or updates existent shortcut.",
			Usage:       "shortcuts add shortcut -> command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				ss := strings.Split(q, "->")
				if len(ss) < 2 {
					_, err = sg.RespondBadCommandUsage(m, c)
					return
				}
				shortcut := strings.TrimSpace(ss[0])
				commandQ := strings.TrimSpace(ss[1])

				// Try to find command.
				command, err := sg.rootCommand.search(sg, commandQ, m)
				if command == nil {
					sg.respondCommandNotFound(m)
					return
				}

				sg.shortcuts.set(sg, shortcut, commandQ)
				sg.RespondTextMention(m, "\""+shortcut+"\" is now a shortcut for \""+commandQ+"\".")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified shortcut.",
			Usage:       "shortcut",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				_, exists := sg.shortcuts.get(sg, q)

				if !exists {
					sg.RespondTextMention(m, "Shortcut \""+q+"\" not found.")
					return
				}

				sg.shortcuts.del(sg, q)
				sg.RespondTextMention(m, "Shortcut \""+q+"\" was deleted successfully.")
				return
			},
		},
	},
}
