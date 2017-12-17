package test

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
)

// Test is just a testing command
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:             "test",
		PermittedByDefault:  true,
		Description:         "Test command.",
		TextResponse:        "Test passed.",
		AllowDefaultChannel: true,
		SubCommands: []*sugo.Command{
			{
				Trigger:            "test1",
				PermittedByDefault: true,
				Description:        "subTest1 command.",
				TextResponse:       "subTest1 passed.",
				SubCommands: []*sugo.Command{
					{
						Trigger:            "test11",
						PermittedByDefault: true,
						Description:        "subTest11 command.",
						TextResponse:       "subTest11 passed.",
					},
					{
						Trigger:            "test12",
						PermittedByDefault: true,
						Description:        "subTest12 command.",
						TextResponse:       "subTest12 passed.",
					},
				},
			},
			{
				Trigger:            "test2",
				PermittedByDefault: true,
				ParamsAllowed:      true,
				Description:        "subTest2 command.",
				TextResponse:       "subTest2 passed.",
			},
			{
				Trigger:            "responses",
				PermittedByDefault: true,
				Description:        "subTest2 command.",
				Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
					if _, err := sg.Respond(m, "default"); err != nil {return err}
					if _, err := sg.RespondInfo(m, "info"); err != nil {return err}
					if _, err := sg.RespondSuccess(m, "success"); err != nil {return err}
					if _, err := sg.RespondWarning(m, "warning"); err != nil {return err}
					if _, err := sg.RespondDanger(m, "failure"); err != nil {return err}
					return nil
				},
			},
		},
	},
}
