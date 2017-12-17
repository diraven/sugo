package feeds

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
)

var rootCommand = &sugo.Command{
	Trigger:     "feeds",
	Description: "Allows to manipulate feeds.",
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		embed := &discordgo.MessageEmbed{
			Title:  "Feeds",
			Fields: []*discordgo.MessageEmbedField{},
		}

		for _, item := range feeds {
			embed.Description = embed.Description + item.Url + "\n"
		}

		_, err = sg.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:       "add",
			Description:   "Adds new feed.",
			Usage:         "http://example.com/rss/",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				// Validate feed url.
				_, err = fp.ParseURL(q)
				if err != nil {
					_, err = sg.RespondDanger(m, "", err.Error())
					return err
				}

				// Add requested URL to the list.
				feeds.add(sg, m.ChannelID, q)

				_, err = sg.RespondSuccess(m, "", "")
				return err
			},
		},
		{
			Trigger:       "del",
			Description:   "Deletes specified feed.",
			Usage:         "http://example.com/rss",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				if err = feeds.del(sg, m.ChannelID, q); err != nil {
					return err
				}

				if _, err = sg.RespondSuccess(m, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
