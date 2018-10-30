package sugo

import (
	"github.com/pkg/errors"
	"log"
)

// HandleError handles unexpected errors that were returned unhandled elsewhere.
func (sg *Instance) HandleError(
	err error,
) {
	// If there is custom error handler:
	if sg.ErrorHandler != nil {
		// Run it.
		sg.ErrorHandler(err)
		return
	}

	// Check if error is sugo error.
	var cause = errors.Cause(err)

	if sugoErr, ok := cause.(*Error); ok {
		// Error is sugo error.
		// Make sure we have got a valid Request to work with.
		if sugoErr.request != nil {
			// Now try to respond with the Embed.
			if _, channelSendErr := sugoErr.request.NewResponse(ResponseDanger, "", sugoErr.text).Send(); channelSendErr != nil {
				// If we were unable to send the message to the same channel command was issued on,
				// try to send to the user DM instead.
				if _, dmSendErr := sugoErr.request.NewResponse(ResponseDanger, "", sugoErr.text).SendDM(); dmSendErr != nil {
					// We were unable to send the error via DM either. In that case just log it into the console as well
					// as all the rest of the errors we have encountered.
					log.Println(channelSendErr)
					log.Println(dmSendErr)
					log.Println(sugoErr)
				}
				// Message was sent to the DM. There is nothing else we need to do.
				return
			}
			// Message was sent to the channel. There is nothing else we need to do.
			return
		}
		// There is a sugoErr, but without Request provided. There must be something very wrong. Report the issue.
		log.Println(sugoErr)
		log.Println(errors.New("sugoErr is provided without Request"))
		return
	}

	// Otherwise just put error into the log.
	log.Println(err)
}
