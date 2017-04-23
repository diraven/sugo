package sugo

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

var Bot Instance

const VERSION string = "0.0.18"

const PermissionNone = 0 // A permission that is always granted.

type Instance struct {
	*discordgo.Session
	Self     *discordgo.User
	root     *discordgo.User
	commands []Command
	Data     data
	Done chan bool
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

func (sg *Instance) RegisterCommand(command Command) (err error) {
	sg.commands = append(sg.commands, command)
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

	// Dispatch command.
	for _, command := range Bot.commands {
		// Test if command is applicable.
		is_applicable, err := command.IsApplicable(&Bot, m.Message)
		if err != nil {
			// TODO: Report error.
		}
		if is_applicable {
			// Check if user has all necessary permissions.
			is_allowed, err := command.IsAllowed(&Bot, m.Message)
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
}
