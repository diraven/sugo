package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
)

// HandleError handles unexpected errors that were returned unhandled elsewhere.
func (sg *Instance) HandleError(req *Request, err error) {
	// If there is custom error handler:
	if sg.ErrorHandler != nil {
		// Run it.
		sg.ErrorHandler(req, err)
		return
	}

	// If request is available:
	if req != nil {
		// Analyze error and check if it's permission issue.

		// First of all, get the underlying error if any.
		var cause error
		cause = errors.Cause(err)
		if cause != nil {
			// Check if error cause is discordgo error.
			if dgoError, ok := cause.(*discordgo.RESTError); ok {
				// Check if error is permissions error.
				if dgoError.Message.Code == 50013 {
					err = errors.Wrap(err, "**bot is missing necessary permissions, contact server admin or responsible person to fix this**")
					// It's a permission error. Try to send PM to the user explaining the issue.
					// Now try to respond with the Embed.
					if _, channelSendErr := req.NewResponse(ResponseDanger, "", err.Error()).Send(); channelSendErr != nil {
						// If we were unable to send the message to the same channel command was issued on,
						// try to send to the user DM instead.
						if _, dmSendErr := req.NewResponse(ResponseDanger, "", err.Error()).SendDM(); dmSendErr != nil {
							// We were unable to send the error via DM either. In that case just log it into the console as well
							// as all the rest of the errors we have encountered.
							log.Println(channelSendErr.Error())
							log.Println(dmSendErr.Error())
							log.Println(err.Error())
						}
						// Message was sent to the DM. There is nothing else we need to do.
						return
					}
					// Message was sent to the channel. There is nothing else we need to do.
					return
				}
			}
		}
	}

	// Otherwise just put error into the log.
	log.Println(err)
}
