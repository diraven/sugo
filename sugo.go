package sugo

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

var Bot Instance

const VERSION string = "0.0.12"

type Instance struct {
	*discordgo.Session
	Self     *discordgo.User
	commands []Command
}

func Start(token string) (sg *Instance) {
	// Create empty Instance session.
	Bot = Instance{}
	sg = &Bot

	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Save Discord session into Instance struct.
	sg.Session = s

	// Get the account information.
	self, err := s.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Save reference to the bot into Instance struct.
	sg.Self = self

	// Register messageCreate as a callback for the messageCreate events.
	s.AddHandler(onMessageCreate)

	// Open the websocket and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	return
}

func (sg *Instance) RegisterCommand(command Command) (err error) {
	sg.commands = append(sg.commands, command)
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	// Simple dispatch commmand.
	for _, command := range Bot.commands {
		if command.Test(&Bot, m.Message) {
			command.Execute(&Bot, m.Message)
		}
	}
}
