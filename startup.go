package sugo

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Startup starts the bot up.
func (sg *Instance) Startup(token string) error {
	// Intitialize Shutdown channel.
	sg.done = make(chan os.Signal, 1)

	// Variable to store errors.
	var err error

	// Create a new Discord Session using the provided bot token.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return errors.New("Error creating Discord Session... " + err.Error())
	}

	// Save Discord Session into Instance struct.
	sg.Session = s

	// Get bot discordgo.User instance.
	self, err := sg.Session.User("@me")
	if err != nil {
		return errors.New("Error obtaining bot account details... " + err.Error())
	}
	sg.Self = self

	// Run startup handlers.
	for _, handler := range sg.startupHandlers {
		if err = handler(sg); err != nil {
			// If there is any error - we stop the startup process and shut the bot down as there is not much sense
			// to let bot finish the startup in an event of an error.
			sg.Shutdown()
			return err
		}
	}

	// Register callback for the messageCreate events.
	sg.Session.AddHandler(func(s *discordgo.Session, mc *discordgo.MessageCreate) {
		sg.onMessageCreate(mc.Message)
	})

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
	sg.shutdown()

	return nil
}
