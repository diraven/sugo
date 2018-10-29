package test

import (
	"github.com/diraven/sugo"
)

// Init initializes module on the given bot.
func Init(sg *sugo.Instance) (err error) {
	return sg.AddCommand(cmd)
}

var cmd = &sugo.Command{
	Trigger:     "test",
	Description: "command made for testing purposes\n",
	Execute: func(req *sugo.Request) (err error) {
		if _, err := req.Respond("", sugo.NewInfoEmbed(req, "test"), false); err != nil {
			return sugo.WrapError(req, err)
		}
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger: "subtest1",
			Execute: func(req *sugo.Request) (err error) {
				if _, err := req.Respond("", sugo.NewInfoEmbed(req, "subtest1"), false); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
			SubCommands: []*sugo.Command{
				{
					Trigger: "subtest11",
					Execute: func(req *sugo.Request) (err error) {
						if _, err := req.Respond("", sugo.NewInfoEmbed(req, "subtest11"), false); err != nil {
							return sugo.WrapError(req, err)
						}
						return
					},
				},
				{
					Trigger: "subtest12",
					Execute: func(req *sugo.Request) (err error) {
						if _, err := req.Respond("", sugo.NewInfoEmbed(req, "subtest12"), false); err != nil {
							return sugo.WrapError(req, err)
						}
						return
					},
				},
			},
		},
		{
			Trigger: "subtest2",
			Execute: func(req *sugo.Request) (err error) {
				if _, err := req.Respond("", sugo.NewInfoEmbed(req, "subtest2"), false); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "responses",
			Execute: func(req *sugo.Request) (err error) {
				if _, err := req.Respond("text", nil, false); err != nil {
					return sugo.WrapError(req, err)
				}
				if _, err := req.Respond("", sugo.NewDefaultEmbed(req, "default"), false); err != nil {
					return sugo.WrapError(req, err)
				}
				if _, err := req.Respond("", sugo.NewInfoEmbed(req, "info"), false); err != nil {
					return sugo.WrapError(req, err)
				}
				if _, err := req.Respond("", sugo.NewWarningEmbed(req, "warning"), false); err != nil {
					return sugo.WrapError(req, err)
				}
				if _, err := req.Respond("", sugo.NewDangerEmbed(req, "danger"), false); err != nil {
					return sugo.WrapError(req, err)
				}
				if _, err := req.Respond("", sugo.NewDefaultEmbed(req, "default (DM)"), true); err != nil {
					return sugo.WrapError(req, err)
				}
				return sugo.NewError(req, "error")
			},
		},
	},
}
