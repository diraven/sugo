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
	all() (shortcuts []sShortcut)
	add(sg *Instance, short string, long string)
	get(sg *Instance, i int) (shortcut *sShortcut, exists bool)
	del(sg *Instance, i int)
	swap(sg *Instance, i1 int, i2 int)
	load(sg *Instance) (data_length int, err error)
	save(sg *Instance) (data_length int, err error)
	startup(sg *Instance) error
	teardown(sg *Instance) error
}

const SHORTCUTS_DATA_FILENAME = "shortcuts.json"

type sShortcut struct {
	Short string
	Long  string
}

type sShortcutsStorage struct {
	Shortcuts []sShortcut
}

func (p *sShortcutsStorage) all() (shortcuts []sShortcut) {
	return p.Shortcuts
}

func (p *sShortcutsStorage) swap(sg *Instance, i1 int, i2 int) {
	p.Shortcuts[i1], p.Shortcuts[i2] = p.Shortcuts[i2], p.Shortcuts[i1]
}

func (p *sShortcutsStorage) get(sg *Instance, i int) (shortcut *sShortcut, exists bool) {
	if i < 0 || i >= len(p.Shortcuts) {
		return nil, false
	}
	return &p.Shortcuts[i], true
}

func (p *sShortcutsStorage) add(sg *Instance, short string, long string) {
	p.Shortcuts = append(p.Shortcuts, sShortcut{
		short,
		long,
	})
}

func (p *sShortcutsStorage) del(sg *Instance, i int) {
	p.Shortcuts = append(p.Shortcuts[:i], p.Shortcuts[i+1:]...)
}

func (p *sShortcutsStorage) load(sg *Instance) (data_length int, err error) {
	if _, error_type := os.Stat(SHORTCUTS_DATA_FILENAME); os.IsNotExist(error_type) {
				// File to load data from does not exist.
		// Check if perms storage is empty and initialize it.
		shortcutsStorage := sg.shortcuts.(*sShortcutsStorage)
		if shortcutsStorage.Shortcuts == nil {
			shortcutsStorage.Shortcuts = make([]sShortcut, 0)
		}
		return
	}

	// Load file.
	data, err := ioutil.ReadFile(SHORTCUTS_DATA_FILENAME)
	if err != nil {
		return
	}

	// Decode JSON data.
	json.Unmarshal(data, sg.shortcuts.(*sShortcutsStorage))
	if err != nil {
		return
	}

	// Log the operation results.
	data_length = len(data)

	return
}

func (p *sShortcutsStorage) save(sg *Instance) (data_length int, err error) {
	// Encode our data into JSON.
	data, err := json.Marshal(sg.shortcuts.(*sShortcutsStorage))
	if err != nil {
		return
	}

	// Save data into file.
	err = ioutil.WriteFile(SHORTCUTS_DATA_FILENAME, data, 0644)
	if err != nil {
		return
	}

	data_length = len(data)

	return
}

func (p *sShortcutsStorage) startup(sg *Instance) (err error) {
	_, err = p.load(sg)
	return
}

func (p *sShortcutsStorage) teardown(sg *Instance) (err error) {
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

		for i, shortcut := range sg.shortcuts.(*sShortcutsStorage).Shortcuts {
			result = result + strconv.FormatInt(int64(i), 10) + ": `" + shortcut.Short + "` -> `" + shortcut.Long + "`\n"
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
			Trigger:     "add",
			Description: "Adds new or updates existent shortcut.",
			Usage:       "shortcuts add shortcut -> command [subcommand ...]",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				ss := strings.Split(q, "->")
				if len(ss) < 2 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}
				short := strings.TrimSpace(ss[0])
				long := strings.TrimSpace(ss[1])

				// Try to find command.
				command, err := sg.rootCommand.search(sg, long, m)
				if command == nil {
					sg.respondCommandNotFound(m)
					return
				}

				sg.shortcuts.add(sg, short, long)
				sg.RespondTextMention(m, "\""+short+"\" is now a shortcut for \""+long+"\".")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified shortcut.",
			Usage:       "1",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				i, err := strconv.ParseInt(q, 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}
				_, exists := sg.shortcuts.get(sg, int(i))

				if !exists {
					sg.RespondTextMention(m, "Shortcut \""+q+"\" not found.")
					return
				}

				sg.shortcuts.del(sg, int(i))
				_, err = sg.RespondSuccessMention(m, "")
				return
			},
		},
		{
			Trigger:     "swap",
			Description: "Swaps specified shortcuts.",
			Usage:       "1 2",
			Execute: func(ctx context.Context, c *Command, q string, sg *Instance, m *discordgo.Message) (err error) {
				var exists bool
				ss := strings.Split(q, " ")
				if len(ss) < 2 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				i1, err := strconv.ParseInt(ss[0], 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}
				_, exists = sg.shortcuts.get(sg, int(i1))
				if !exists {
					sg.RespondTextMention(m, "Shortcut \""+q+"\" not found.")
					return
				}

				i2, err := strconv.ParseInt(ss[1], 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}
				_, exists = sg.shortcuts.get(sg, int(i2))
				if !exists {
					sg.RespondTextMention(m, "Shortcut \""+q+"\" not found.")
					return
				}

				sg.shortcuts.swap(sg, int(i1), int(i2))

				_, err = sg.RespondSuccessMention(m, "")
				return
			},
		},
	},
}
