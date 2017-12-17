package sugo

import (
	"strings"
	"github.com/bwmarrin/discordgo"
	"context"
	"errors"
)

// onMessageCreate contains all the message processing logic for the bot.
func onPresenceUpdate(s *discordgo.Session, pu *discordgo.PresenceUpdate) {
	for _, module := range Bot.Modules {
		if module.OnPresenceUpdate != nil {
			if err := module.OnPresenceUpdate(Bot, pu); err != nil {
				Bot.HandleError(err)
			}
		}
	}
}
