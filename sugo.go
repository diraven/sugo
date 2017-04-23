package sugo

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"github.com/diraven/sugo/errors"
	"strings"
	"github.com/diraven/sugo/helpers"
)

const VERSION string = "0.0.20"

const PermissionNone = 0 // A permission that is always granted.

var Bot Instance

type Instance struct {
	*discordgo.Session
	Self     *discordgo.User
	root     *discordgo.User
	commands map[string]Command
	data     *bot_data
	Done     chan bool
}

func (sg *Instance) Startup(token string, root_uid string) (err error) {
	// Intitialize Done channel.
	sg.Done = make(chan bool)

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

	// Get bot account info.
	self, err := sg.Session.User("@me")
	if err != nil {
		fmt.Println("Error obtaining account details... ", err)
		return
	}
	sg.Self = self

	// Get root account info.
	if root_uid != "" {
		root, err := sg.Session.User(root_uid)
		if err != nil {
			// TODO: Report error.
		} else {
			sg.root = root
		}
	}

	// Register callback for the messageCreate events.
	sg.Session.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = sg.Session.Open()
	if err != nil {
		fmt.Println("Error opening connection... ", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Block until a message from Done channel.
	<-sg.Done
	return
}

func (sg *Instance) Shutdown() (err error) {
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

	// Send to Done channel.
	sg.Done <- true
	return
}

func (sg *Instance) RegisterCommand(trigger string, c Command) (err error) {
	// Initialize commands storage.
	if sg.commands == nil {
		sg.commands = make(map[string]Command)
	}

	if _, ok := sg.commands[trigger]; ok {
		return errors.SugoError{
			Text: fmt.Sprintf("Conflicting triggers: Command with top level '%s' trigger already exists.", trigger),
		}
	}
	sg.commands[trigger] = c
	return
}

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

func (sg *Instance) UserHasPermission(permission int, u *discordgo.User, c *discordgo.Channel) (result bool, err error) {
	perms, err := sg.UserChannelPermissions(u.ID, c.ID)
	if err != nil {
		return
	}
	result = (perms | permission) == perms
	return
}

func (sg *Instance) BotHasPermission(permission int, c *discordgo.Channel) (result bool, err error) {
	result, err = sg.UserHasPermission(permission, sg.Self, c)
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		// TODO: Report error.
		return
	}

	// Make sure message is not sent by bot.
	if m.Author.Bot {
		return
	}

	// Make sure the bot is mentioned in the message, and bot mention is first mention in the message.
	if len(m.Mentions) < 1 {
		return
	}
	if m.Mentions[0].ID != Bot.Self.ID {
		return
	}

	// Consume bot mention from the message content.
	m.Content = strings.TrimSpace(m.Content)
	m.Content = strings.TrimPrefix(m.Content, fmt.Sprintf("%s", helpers.UserAsMention(Bot.Self)))
	m.Content = strings.TrimSpace(m.Content)

	// Try to figure out command name.
	next_space_index := strings.Index(m.Content, " ")
	var command_name string
	if next_space_index < 0 {
		command_name = m.Content
	} else {
		command_name = m.Content[:strings.Index(m.Content, " ")]
	}

	// Consume command name.
	m.Content = strings.TrimPrefix(m.Content, fmt.Sprintf("%s", command_name))
	m.Content = strings.TrimSpace(m.Content)

	// Dispatch command.
	command, ok := Bot.commands[command_name]
	if ok {
		is_allowed, err := command.CheckPermissions(&Bot, m.Message)
		if err != nil {
			// TODO: Report error.
		}
		if is_allowed {
			// Execute command.
			err := command.Execute(&Bot, m.Message)
			if err != nil {
				// TODO: Report error.
			}
		}
	}
}
