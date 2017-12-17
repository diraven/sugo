package greet

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

// TODO Fix the greeting module.

// Greet responds to the user with greeting and invitation to use help command.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "", // Command with no trigger will be applied if message consists from bot mention only.
		PermittedByDefault: true,
		Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
			var err error

			_, err = sg.RespondInfo(m, "", "Hi! My name is "+
				fmt.Sprintf("%s and I'm here to help you out... ", sg.Self.Username)+
				"Try `help` for more info.",
			)
			return err
		},
	},
}
