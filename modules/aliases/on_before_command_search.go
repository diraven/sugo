package aliases

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strings"
)

func onBeforeCommandSearch(sg *sugo.Instance, req *sugo.Request) error {
	// We only work with guild text channels and ignore everything else.
	if req.Channel.Type != discordgo.ChannelTypeGuildText {
		return nil
	}

	// Process aliases.
	for alias, commandPath := range *aliases.all(req.Guild) {
		if strings.Index(req.Query, alias) == 0 {
			if len(req.Query) == len(alias) || string(req.Query[len(alias)]) == " " {
				req.Query = strings.Replace(req.Query, alias, commandPath, 1)
				break
			}
		}
	}
	// Return resulting query.
	return nil
}
