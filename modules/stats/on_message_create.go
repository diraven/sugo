package stats

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(sg *sugo.Instance, mc *discordgo.MessageCreate) error {
	guild, err := sg.GuildFromMessage(mc.Message)
	if err != nil {
		return err
	}

	// Ignore bots.
	if mc.Author.Bot {
		return nil
	}

	stats.logMessage(sg, guild.ID, mc.Author.ID)

	return nil
}
