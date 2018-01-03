package aliases

import (
	"github.com/diraven/sugo"
	"strings"
	"github.com/bwmarrin/discordgo"
)

func onBeforeCommandSearch(sg *sugo.Instance, m *discordgo.Message, q string) (string, error) {
	// Get channel.
	channel, err := sg.ChannelFromMessage(m)
	if err != nil {
		return "", err
	}

	// We only work with guild text channels and ignore everything else.
	if channel.Type != discordgo.ChannelTypeGuildText {
		return q, nil
	}

	// Get guild.
	guild, err := sg.GuildFromMessage(m)
	if err != nil {
		return "", err
	}

	// Process aliases.
	for alias, commandPath := range *aliases.all(guild) {
		if strings.Index(q, alias) == 0 {
			if len(q) == len(alias) || string(q[len(alias)]) == " " {
				q = strings.Replace(q, alias, commandPath, 1)
				break
			}
		}
	}
	// Return resulting query.
	return q, nil
}
