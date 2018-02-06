package sugo

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"github.com/pkg/errors"
)

func (sg *Instance) onMessageCreate(m *discordgo.Message) {
	var err error // Used to capture and report errors.

	// Instance request.
	var r = &Request{}

	// Put bot pointer into the appropriate var for reference.
	r.Sugo = sg

	// Put message into request.
	r.Message = m

	// Put initial query into request.
	r.Query = m.Content

	// Ignore any message that is coming from bot.
	if m.Author.Bot {
		return
	}

	// Get channel.
	r.Channel, err = r.Sugo.Session.State.Channel(r.Message.ChannelID)
	if err != nil {
		sg.HandleError(errors.Wrap(err, "getting channel failed"))
	}

	if r.Channel.Type == discordgo.ChannelTypeDM {
		// It's Direct Messaging Channel. Every message here is in fact a direct message to the bot, so we consider
		// it to be command without any further checks.

	} else if r.Channel.Type == discordgo.ChannelTypeGuildText || r.Channel.Type == discordgo.ChannelTypeGroupDM {
		// It's either Guild Text Channel or multiple people direct group Channel.
		// In order to detect command we need to account for Trigger.

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

		// If the message starts with bot mention.
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
		// Unhandled error in command.
		sg.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + r.Query + ")"))
	}

	// If we did not find matching command, try applying alias and searching again.
	if r.Command == nil {
		// Apply aliases if any applicable.
		for _, alias := range *sg.aliases {
			if strings.HasPrefix(strings.TrimSpace(r.Query), alias.from) {
				r.Query = strings.Replace(r.Query, alias.from, alias.to, 1)
				break // we apply only one alias that matched first.
			}
		}

		// Search for applicable command again after alias applied.
		r.Command, err = sg.FindCommand(r, r.Query)
		if err != nil {
			// Unhandled error in command.
			sg.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + r.Query + ")"))
		}
	}

	if r.Command != nil {
		// Remove command Trigger from message string.
		r.Query = strings.TrimSpace(strings.TrimPrefix(r.Query, r.Command.GetPath()))

		// And execute command.
		err = r.Command.execute(sg, r)
		if err != nil {
			if strings.Contains(err.Error(), "\"code\": 50013") {
				// Insufficient permissions, bot configuration issue.
				sg.HandleError(errors.New("Bot permissions error: " + err.Error() + " (" + r.Query + ")"))
			} else {
				// Other discord errors.
				sg.HandleError(errors.New("Bot command execute error: " + err.Error() + " (" + r.Query + ")"))
			}
			sg.HandleError(err)
		}
	}

	// Command not found.
}
