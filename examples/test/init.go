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
	Execute: func(sg *sugo.Instance, r *sugo.Request) error {
		_, err := sg.RespondInfo(r, "", "test passed")
		return err
	},
	SubCommands: []*sugo.Command{
		{
			Trigger: "subtest1",
			Execute: func(sg *sugo.Instance, r *sugo.Request) error {
				_, err := sg.RespondInfo(r, "", "subtest1 passed")
				return err
			},
			SubCommands: []*sugo.Command{
				{
					Trigger: "subtest11",
					Execute: func(sg *sugo.Instance, r *sugo.Request) error {
						_, err := sg.RespondInfo(r, "", "subtest11 passed")
						return err
					},
				},
				{
					Trigger: "subtest12",
					Execute: func(sg *sugo.Instance, r *sugo.Request) error {
						_, err := sg.RespondInfo(r, "", "subtest12 passed")
						return err
					},
				},
			},
		},
		{
			Trigger: "subtest2",
			Execute: func(sg *sugo.Instance, r *sugo.Request) error {
				_, err := sg.RespondInfo(r, "", "subtest2 passed")
				return err
			},
		},
		{
			Trigger: "responses",
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
}
