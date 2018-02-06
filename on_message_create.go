package sugo

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"errors"
)

// onMessageCreate is a lowest level handler for bot. All the request building and command searching magic happen here.
func (sg *Instance) onMessageCreate(m *discordgo.Message) {
	var err error

	// Ignore a message it's author is bot.
	if m.Author.Bot {
		return
	}

	// Instantiate request.
	var r = &Request{}

	// Put bot pointer into the appropriate request var for later reference.
	r.Sugo = sg

	// Put message into request.
	r.Message = m

	// Put initial query into request.
	r.Query = m.Content

	// Get message channel and put it into the request.
	r.Channel, err = sg.Session.State.Channel(r.Message.ChannelID)
	if err != nil {
		sg.HandleError(errors.New("unable to get channel: " + err.Error() + " (" + r.Query + ")"))
	}

	if r.Channel.Type == discordgo.ChannelTypeDM {
		// It's Direct Messaging Channel. Every message here is in fact a direct message to the bot, so we consider
		// it to be command without any further checks for prefixes.

	} else if r.Channel.Type == discordgo.ChannelTypeGuildText || r.Channel.Type == discordgo.ChannelTypeGroupDM {
		// It's either Guild Text Channel or multiple people direct group Channel.
		// In order to detect command we need to check for bot Trigger.

		// If bot Trigger is set and command starts with that Trigger:
		if sg.Trigger != "" && strings.HasPrefix(r.Query, sg.Trigger) {
			// Replace custom Trigger with bot mention for it to be detected as bot Trigger.
			r.Query = strings.Replace(r.Query, sg.Trigger, sg.Self.Mention(), 1)
		}

		// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
		// for mention detection to work right.
		if strings.HasPrefix(r.Query, "<@!") {
			r.Query = strings.Replace(r.Query, "<@!", "<@", 1)
		}

		// If the message starts with bot mention:
		if strings.HasPrefix(strings.TrimSpace(r.Query), sg.Self.Mention()) {
			// Remove bot Trigger from the string.
			r.Query = strings.TrimSpace(strings.TrimPrefix(r.Query, sg.Self.Mention()))
		} else {
			// Ignore the message otherwise and do nothing.
			return
		}

	}

	// Search for applicable command.
	r.Command, err = sg.FindCommand(r, r.Query)
	if err != nil {
		sg.HandleError(errors.New("command search error: " + err.Error() + " (" + r.Query + ")"))
	}

	// If we did not find matching command, try applying alias and search again.
	if r.Command == nil {
		// Apply aliases if any applicable.
		for _, alias := range *sg.aliases {
			if strings.HasPrefix(strings.TrimSpace(r.Query), alias.from) {
				r.Query = strings.Replace(r.Query, alias.from, alias.to, 1)
				break // we apply only one alias that matched first.
			}
		}

		// Search for applicable command again after alias was applied.
		r.Command, err = sg.FindCommand(r, r.Query)
		if err != nil {
			sg.HandleError(errors.New("command search error: " + err.Error() + " (" + r.Query + ")"))
		}
	}

	// If we have found applicable command:
	if r.Command != nil {
		// Remove command Trigger from message string.
		r.Query = strings.TrimSpace(strings.TrimPrefix(r.Query, r.Command.GetPath()))

		// And execute command.
		err = r.Command.execute(sg, r)
		if err != nil {
			if strings.Contains(err.Error(), "\"code\": 50013") {
				// Insufficient permissions.
				sg.HandleError(errors.New("permissions error: " + err.Error() + " (" + r.Query + ")"))
			}
			sg.HandleError(errors.New("command execution error: " + err.Error() + " (" + r.Query + ")"))
		}
	}

	// Command not found, we do nothing.
}
