package test

import (
	"github.com/diraven/sugo"
)

// Init initializes module on the given bot.
func Init(sg *sugo.Instance) {
	sg.AddCommand(cmd)
}

var cmd = &sugo.Command{
	Trigger:     "test",
	Description: "command made for testing purposes\n asdf fdsa",
	Execute: func(req *sugo.Request) error {
		_, err := req.RespondInfo("", "test passed")
		return err
	},
	SubCommands: []*sugo.Command{
		{
			Trigger: "subtest1",
			Execute: func(req *sugo.Request) error {
				_, err := req.RespondInfo("", "subtest1 passed")
				return err
			},
			SubCommands: []*sugo.Command{
				{
					Trigger: "subtest11",
					Execute: func(req *sugo.Request) error {
						_, err := req.RespondInfo("", "subtest11 passed")
						return err
					},
				},
				{
					Trigger: "subtest12",
					Execute: func(req *sugo.Request) error {
						_, err := req.RespondInfo("", "subtest12 passed")
						return err
					},
				},
			},
		},
		{
			Trigger: "subtest2",
			Execute: func(req *sugo.Request) error {
				_, err := req.RespondInfo("", "subtest2 passed")
				return err
			},
		},
		{
			Trigger: "responses",
			Execute: func(req *sugo.Request) error {
				if _, err := req.Respond("", "default", sugo.ColorPrimary, ""); err != nil {
					return err
				}
				if _, err := req.RespondInfo("", "info"); err != nil {
					return err
				}
				if _, err := req.RespondSuccess("", "success"); err != nil {
					return err
				}
				if _, err := req.RespondWarning("", "warning"); err != nil {
					return err
				}
				if _, err := req.RespondDanger("", "failure"); err != nil {
					return err
				}
				return nil
			},
		},
	},
}
