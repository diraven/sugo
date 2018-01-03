package stats

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(sg *sugo.Instance, m *discordgo.Message) error {
	// Get channel.
	channel, err := sg.ChannelFromMessage(m)
	if err != nil {
		return err
	}

	// We only work with guild text channels and ignore everything else.
	if channel.Type != discordgo.ChannelTypeGuildText {
		return nil
	}

	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return err
	}

	// Ignore bots.
	if m.Author.Bot {
		return nil
	}

	stats.logMessage(sg, guild.ID, m.Author.ID)

	return nil
}
