package sugo

import (
	"log"
	"os/signal"
	"syscall"
	"os"
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"errors"
)

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

	// Register callback for the messageUpdate events.
	//sg.Session.AddHandler(onMessageUpdate)

	// Register callback for the presenceUpdate events.
	sg.Session.AddHandler(onPresenceUpdate)

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

	return nil
}
