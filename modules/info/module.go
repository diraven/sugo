package info

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
	"github.com/dustin/go-humanize"
	"runtime"
	"time"
)

// Info shows some general bot info.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:             "info",
		PermittedByDefault:  true,
		AllowDefaultChannel: true,
		Description:         "Shows basic info about bot.",
		Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
			var err error

			// Set command response.
			now := time.Now().UTC()
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			_, err = sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				URL:         "https://github.com/diraven/sugo",
				Title:       "https://github.com/diraven/sugo",
				Description: "A Discord bot written in Go.",
				Timestamp:   helpers.TimeToDiscordTimestamp(now),
				Color:       sugo.ColorInfo,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "DiRaven#0519",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "sugo ver.",
						Value:  sugo.VERSION,
						Inline: true,
					},
					{
						Name:   "discordgo ver.",
						Value:  discordgo.VERSION,
						Inline: true,
					},
					{
						Name:   "go ver.",
						Value:  runtime.Version(),
						Inline: true,
					},
					{
						Name:   "memory usage",
						Value:  fmt.Sprintf("%s", humanize.Bytes(memStats.Sys)),
						Inline: true,
					},
					{
						Name:   "goroutines",
						Value:  fmt.Sprintf("%d", runtime.NumGoroutine()),
						Inline: true,
					},
					{
						Name:   "guilds",
						Value:  fmt.Sprintf("%d", len(sg.State.Guilds)),
						Inline: true,
					},
					{
						Name:   "modules",
						Value:  fmt.Sprintf("%d", len(sg.Modules)),
						Inline: true,
					},
				},
			})
			if err != nil {
				return err
			}
			return nil
		},
	},
}
