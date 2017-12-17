package stats

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

func onPresenceUpdate(sg *sugo.Instance, pu *discordgo.PresenceUpdate) error {
	if pu.Game != nil {
		if err := stats.logPlaying(sg, pu.GuildID, pu.User.ID, pu.Game.Type, pu.Game.Name); err != nil {
			return err
		}
	}

	return nil
}
