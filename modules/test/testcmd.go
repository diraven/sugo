package test

import (
	"github.com/diraven/sugo"
)

// Module is just a few testing commands
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
				AllowParams:        true,
				Description:        "subTest2 command.",
				TextResponse:       "subTest2 passed.",
			},
			{
				Trigger:            "responses",
				PermittedByDefault: true,
				Description:        "subTest2 command.",
				Execute: func(sg *sugo.Instance, req *sugo.Request) error {
					if _, err := sg.Respond(req, "", "default", sugo.ColorPrimary, ""); err != nil {
						return err
					}
					if _, err := sg.RespondInfo(req, "", "info"); err != nil {
						return err
					}
					if _, err := sg.RespondSuccess(req, "", "success"); err != nil {
						return err
					}
					if _, err := sg.RespondWarning(req, "", "warning"); err != nil {
						return err
					}
					if _, err := sg.RespondDanger(req, "", "failure"); err != nil {
						return err
					}
					return nil
				},
			},
		},
	},
}
