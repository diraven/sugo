package stats

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

func onMessageCreate(sg *sugo.Instance, req *sugo.Request) error {
	// We only work with guild text channels and ignore everything else.
	if req.Channel.Type != discordgo.ChannelTypeGuildText {
		return nil
	}

	// Ignore bots.
	if req.Message.Author.Bot {
		return nil
	}

	stats.logMessage(sg, req.Guild.ID, req.Message.Author.ID)

	return nil
}
