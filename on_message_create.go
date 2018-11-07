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

	// Search for applicable command.
	req.Command, err = sg.FindCommand(req, req.Query)
	if err != nil {
		sg.HandleError(req, errors.Wrap(err, "unable to search commands"))
	}

	//// If we did not find matching command, try applying alias and search again.
	//if req.Command == nil {
	//	// Apply aliases if any applicable.
	//	for _, alias := range *sg.aliases {
	//		if strings.HasPrefix(strings.TrimSpace(req.Query), alias.from) {
	//			req.Query = strings.Replace(req.Query, alias.from, alias.to, 1)
	//			break // we apply only one alias that matched first.
	//		}
	//	}
	//
	//	// Search for applicable command again after alias was applied.
	//	req.Command, err = sg.FindCommand(req, req.Query)
	//	if err != nil {
	//		sg.HandleError(errors.Wrap(err, "aliased command search error"), req)
	//	}
	//}

	// If we have found applicable command:
	if req.Command != nil {
		// Remove command Trigger from message string.
		req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Command.GetPath()))

		// And execute command.
		err = req.Command.execute(sg, req)
		if err != nil {
			sg.HandleError(req, errors.Wrap(err, "command execution error"))
		}
	}

	// Command not found, we do nothing.
	return
}
