package sugo

import (
	"os"
	"log"
)

// Shutdown sends Shutdown signal to the bot's Shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before Shutdown.
func (sg *Instance) teardown() error {
	// Perform teardown for all Modules.
	for _, module := range sg.Modules {
		if err := module.teardown(sg); err != nil {
			log.Println(err)
		}
	}

	// Close DB connection.
	sg.DB.Close()

	// Close discord session.
	if err := sg.Session.Close(); err != nil {
		return err
	}
	return nil
}
