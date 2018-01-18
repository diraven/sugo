package sugo

import (
	"github.com/bwmarrin/discordgo"
)

// onMessageUpdate just redirects the call to the onMessageCreate.
func onMessageUpdate(s *discordgo.Session, mu *discordgo.MessageUpdate) {
	var err error
	var message *discordgo.Message

	message, err = s.State.Message(mu.Message.ChannelID, mu.Message.ID)
	if err != nil {
		Bot.HandleError(err)
	}

	Bot.processMessage(message)
}
