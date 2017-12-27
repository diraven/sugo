package sugo

import (
	"github.com/bwmarrin/discordgo"
)

// onMessageCreate contains all the message processing logic for the bot.
func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	Bot.processMessage(mc.Message)
}
