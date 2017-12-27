package stats

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(sg *sugo.Instance, m *discordgo.Message) error {
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
