package sugo

import (
	"github.com/pkg/errors"
	"os"
)

// Shutdown sends Shutdown signal to the bot's Shutdown channel.
func (sg *Instance) Shutdown() {
	sg.done <- os.Interrupt
}

// teardown gracefully releases all resources and saves data before Shutdown.
func (sg *Instance) shutdown() (err error) {
	// Run shutdown handlers.
	for _, handler := range sg.shutdownHandlers {
		if err := handler(sg); err != nil {
			// In case of an error - we report the error and continue the shutdown process. Errors should not interrupt
			// shutdown as we need to perform shutdown as cleanly as possible.
			sg.HandleError(errors.Wrap(err, "shutdown error"))
		}
	}

	// Close discord Session.
	if err := sg.Session.Close(); err != nil {
		return errors.Wrap(err, "discordgo session close error")
	}

	// No errors.
	return
}
