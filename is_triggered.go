package sugo

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (sg *Instance) isTriggered(req *Request) (triggered bool) {
	// Use custom IsTriggered function is provided.
	if sg.IsTriggered != nil {
		return sg.IsTriggered(req)
	}

	if req.Channel.Type == discordgo.ChannelTypeDM {
		// It's Direct Messaging Channel. Every message here is in fact a direct message to the bot, so we consider
		// it to be command without any further checks for prefixes.
		triggered = true
		return
	} else if req.Channel.Type == discordgo.ChannelTypeGuildText || req.Channel.Type == discordgo.ChannelTypeGroupDM {
		// It's either Guild Text Channel or multiple people direct group Channel.
		// In order to detect command we need to check for bot Trigger.

		// If bot Trigger is set and command starts with that Trigger:
		if sg.DefaultTrigger != "" && strings.HasPrefix(req.Query, sg.DefaultTrigger) {
			// Replace custom Trigger with bot mention for it to be detected as bot Trigger.
			req.Query = strings.Replace(req.Query, sg.DefaultTrigger, req.Sugo.Self.Mention(), 1)
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
			triggered = true
			return
		}

		// Otherwise bot is not triggered.
		return
	}

	// We ignore all other channel types and consider bot not triggered.
	return
}
