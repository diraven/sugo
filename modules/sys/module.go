package sys

import (
	"github.com/diraven/sugo"
)

// Module contains system-related commands for bot.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:             "sys",
		Description:         "A set of system-related commands.",
		RootOnly:            true,
		PermittedByDefault:  true,
		AllowDefaultChannel: true,
		SubCommands: []*sugo.Command{
			{
				Trigger:            "shutdown",
				Description:        "Shuts the bot town.",
				PermittedByDefault: true,
				TextResponse:       "Until next time, master!",
				Execute: func(sg *sugo.Instance, req *sugo.Request) error {
					sg.Shutdown()
					return nil
				},
			},
		},
	},
}
