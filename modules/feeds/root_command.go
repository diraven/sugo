package feeds

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var rootCommand = &sugo.Command{
	Trigger:     "feeds",
	Description: "Allows to manipulate feeds.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		embed := &discordgo.MessageEmbed{
			Title:  "Feeds",
			Fields: []*discordgo.MessageEmbedField{},
		}

		for _, item := range feeds {
			embed.Description = embed.Description + item.Url + "\n"
		}

		_, err = sg.ChannelMessageSendEmbed(req.Message.ChannelID, embed)
		return err
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "add",
			Description: "Adds new feed.",
			Usage:       "http://example.com/rss/",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var err error

				// Validate feed url.
				_, err = fp.ParseURL(req.Query)
				if err != nil {
					_, err = sg.RespondDanger(req, "", err.Error())
					return err
				}

				// Add requested URL to the list.
				feeds.add(sg, req.Channel.ID, req.Query)

				_, err = sg.RespondSuccess(req, "", "")
				return err
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified feed.",
			Usage:       "http://example.com/rss",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var err error

				if err = feeds.del(sg, req.Channel.ID, req.Query); err != nil {
					return err
				}

				if _, err = sg.RespondSuccess(req, "", ""); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
