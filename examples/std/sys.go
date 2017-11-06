package std

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

// SYS contains system-related commands for bot.
var SYS = &sugo.Command{
	Trigger:     "sys",
	Description: "A set of system-related commands.",
	RootOnly:    true,
	PermittedByDefault: true,
	AllowDefaultChannel: true,
	SubCommands: []*sugo.Command{
		{
			Trigger:      "shutdown",
			Description:  "Shuts the bot town.",
			PermittedByDefault:  true,
			TextResponse: "Until next time, master!",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				sg.Shutdown()
				return
			},
		},
	},
}
