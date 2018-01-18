package guildwars2

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	url2 "net/url"
	"strings"
)

var infoUrls = map[string]string{
	"GuildWars wiki": "https://wiki.guildwars2.com/index.php?title=Special%%20Search&search=%s&go=Go",
	"GW2Spidy":       "https://www.gw2spidy.com/search/%s?recipes=0",
	"GW2TP":          "https://www.gw2tp.com/search?name=%s",
}

// System returns info about given system.
var itemInfo = &sugo.Command{
	Trigger:            "item_info",
	PermittedByDefault: true,
	Description:        "Provides urls to the websites containing given item info.",
	Usage:              "Item Name",
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		// Make sure there is a query specified.
		if strings.TrimSpace(req.Query) == "" {
			_, err = sg.RespondBadCommandUsage(req, "", "")
			return err
		}

		// Build embed.
		embed := &discordgo.MessageEmbed{
			Title:  "Item info resources",
			Fields: []*discordgo.MessageEmbedField{},
		}

		// Add fields based on urls list.
		for k, v := range infoUrls {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  k,
				Value: fmt.Sprintf(v, url2.QueryEscape(req.Query)),
			})
		}

		// Respond with the embed we just built.
		_, err = sg.ChannelMessageSendEmbed(req.Channel.ID, embed)
		return err
	},
}
