package std

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
var Info = &sugo.Command{
	Trigger:     "info",
	PermittedByDefault: true,
	AllowDefaultChannel: true,
	Description: "Shows basic info about bot.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Set command response.
		now := time.Now().UTC()
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		_, err = sg.RespondEmbed(m, &discordgo.MessageEmbed{
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
					Name:   "sugo",
					Value:  sugo.VERSION,
					Inline: true,
				},
				{
					Name:   "discordgo",
					Value:  discordgo.VERSION,
					Inline: true,
				},
				{
					Name:   "go",
					Value:  runtime.Version(),
					Inline: true,
				},
				{
					Name:   "Memory Usage:",
					Value:  fmt.Sprintf("%s", humanize.Bytes(memStats.Sys)),
					Inline: true,
				},
				{
					Name:   "Goroutines:",
					Value:  fmt.Sprintf("%d", runtime.NumGoroutine()),
					Inline: true,
				},
			},
		})
		if err != nil {
			return
		}
		return
	},
}
