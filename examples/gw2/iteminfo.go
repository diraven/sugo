package gw2

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	url2 "net/url"
	"strings"
)

//fmt.Sprintf("at %v, %s", e.When, e.What)

var infoUrls map[string]string = map[string]string{
	"GuildWars wiki": "https://wiki.guildwars2.com/index.php?title=Special%%20Search&search=%s&go=Go",
	"GW2Spidy":       "https://www.gw2spidy.com/search/%s?recipes=0",
	"GW2TP":          "https://www.gw2tp.com/search?name=%s",
}

// System returns info about given system.
var ItemInfo = &sugo.Command{
	Trigger:            "iteminfo",
	PermittedByDefault: true,
	Description:        "Provides urls to the websites containing given item info.",
	Usage:              "Item Name",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure there is a query specified.
		if strings.TrimSpace(q) == "" {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
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
				Value: fmt.Sprintf(v, url2.QueryEscape(q)),
			})
		}

		// Respond with the embed we just built.
		_, err = sg.RespondEmbed(m, embed)
		return
	},
}