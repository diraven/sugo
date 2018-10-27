package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"strings"
)

func isTriggered(req *Request) bool {
	if req.Channel.Type == discordgo.ChannelTypeDM {
		// It's Direct Messaging Channel. Every message here is in fact a direct message to the bot, so we consider
		// it to be command without any further checks for prefixes.
		return true
	} else if req.Channel.Type == discordgo.ChannelTypeGuildText || req.Channel.Type == discordgo.ChannelTypeGroupDM {
		// It's either Guild Text Channel or multiple people direct group Channel.
		// In order to detect command we need to check for bot Triggereq.

		// If bot Trigger is set and command starts with that Trigger:
		if req.Sugo.Trigger != "" && strings.HasPrefix(req.Query, req.Sugo.Trigger) {
			// Replace custom Trigger with bot mention for it to be detected as bot Triggereq.
			req.Query = strings.Replace(req.Query, req.Sugo.Trigger, req.Sugo.Self.Mention(), 1)
		}

		// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
		// for mention detection to work right.
		if strings.HasPrefix(req.Query, "<@!") {
			req.Query = strings.Replace(req.Query, "<@!", "<@", 1)
		}

		// If the message starts with bot mention:
		if strings.HasPrefix(strings.TrimSpace(req.Query), req.Sugo.Self.Mention()) {
			// Remove bot Trigger from the string.
			req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Sugo.Self.Mention()))
			// Bot is triggered.
			return true
		}

		// Otherwise bot is not triggered.
		return false
	}

	// We ignore all other channel types and consider bot not triggered.
	return false
}

// onMessageCreate is a lowest level handler for bot. All the request building and command searching magic happen here.
func (sg *Instance) onMessageCreate(m *discordgo.Message) {
	var err error

	// Ignore a message it's author is bot.
	if m.Author.Bot {
		return
	}

	// Instantiate request.
	var req = &Request{}

	// Put bot pointer into the appropriate request var for later reference.
	req.Sugo = sg

	// Put message into request.
	req.Message = m

	// Put initial query into request.
	req.Query = m.Content

	// Get message channel and put it into the request.
	req.Channel, err = sg.Session.State.Channel(req.Message.ChannelID)
	if err != nil {
		sg.HandleError(err, req)
	}

	// Make sure bot is triggered by the request.
	if !isTriggered(req) {
		return
	}

	// Search for applicable command.
	req.Command, err = sg.FindCommand(req, req.Query)
	if err != nil {
		sg.HandleError(errors.Wrap(err, "command search error"), req)
	}

	// If we did not find matching command, try applying alias and search again.
	if req.Command == nil {
		// Apply aliases if any applicable.
		for _, alias := range *sg.aliases {
			if strings.HasPrefix(strings.TrimSpace(req.Query), alias.from) {
				req.Query = strings.Replace(req.Query, alias.from, alias.to, 1)
				break // we apply only one alias that matched first.
			}
		}

		// Search for applicable command again after alias was applied.
		req.Command, err = sg.FindCommand(req, req.Query)
		if err != nil {
			sg.HandleError(errors.Wrap(err, "aliased command search error"), req)
		}
	}

	// If we have found applicable command:
	if req.Command != nil {
		// Remove command Trigger from message string.
		req.Query = strings.TrimSpace(strings.TrimPrefix(req.Query, req.Command.GetPath()))

		// And execute command.
		err = req.Command.execute(sg, req)
		if err != nil {
			sg.HandleError(errors.Wrap(err, "command execution error"), req)
		}
	}

	// Command not found, we do nothing.
	return
}
