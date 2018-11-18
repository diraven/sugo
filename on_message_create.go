package sugo

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"strings"
)

// onMessageCreate is a lowest level handler for bot. All the Request building and command searching magic happen here.
func (sg *Instance) onMessageCreate(m *discordgo.Message) {
	var err error

	// Ignore a message it's author is bot.
	if m.Author.Bot {
		return
	}

	// Instantiate Request.
	var req = &Request{}

	// Create Request context.
	req.Ctx = context.Background()

	// Put bot pointer into the appropriate Request var for later reference.
	req.Sugo = sg

	// Put message into Request.
	req.Message = m

	// Put initial query into Request.
	req.Query = m.Content

	// Get message channel and put it into the Request.
	req.Channel, err = sg.Session.State.Channel(req.Message.ChannelID)
	if err != nil {
		sg.HandleError(req, errors.Wrap(err, "unable to retrieve discord channel"))
	}

	// Make sure bot is triggered by the Request.
	if !sg.isTriggered(req) {
		return
	}

	// Apply request middlewares if any.
	for _, m := range sg.requestMiddlewares {
		if err = m(req); err != nil {
			sg.HandleError(req, err)
		}
	}

	// Search for applicable command.
	req.Command, err = sg.FindCommand(req, req.Query)
	if err != nil {
		sg.HandleError(req, errors.Wrap(err, "unable to search commands"))
	}

	// If we have found applicable command:
	if req.Command != nil {
		// Remove command Trigger from message string.
		req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Command.GetPath()))

		// And execute command.
		var resp *Response
		if resp, err = req.Command.execute(sg, req); err != nil {
			sg.HandleError(req, errors.Wrap(err, "command execution error"))
		}

		// Apply response middlewares.
		for _, m := range sg.responseMiddlewares {
			if err = m(resp); err != nil {
				sg.HandleError(req, err)
			}
		}

		// Process response.
		if resp != nil {
			if _, err = resp.Send(); err != nil {
				sg.HandleError(req, errors.Wrap(err, "response processing error"))
			}
		}
	}

	// Command not found, we do nothing.
	return
}
