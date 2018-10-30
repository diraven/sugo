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
		if _, err = req.SimpleResponse("test").Send(); err != nil {
			return sugo.WrapError(req, err)
		}
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger: "subtest1",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.SimpleResponse("subtest1").Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
			SubCommands: []*sugo.Command{
				{
					Trigger: "subtest11",
					Execute: func(req *sugo.Request) (err error) {
						if _, err = req.SimpleResponse("subtest11").Send(); err != nil {
							return sugo.WrapError(req, err)
						}
						return
					},
				},
				{
					Trigger: "subtest12",
					Execute: func(req *sugo.Request) (err error) {
						if _, err = req.SimpleResponse("subtest12").Send(); err != nil {
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
				if _, err = req.SimpleResponse("subtest2").Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "plaintext",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.NewResponse(
					sugo.ResponsePlainText,
					"message title",
					"message text",
				).Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "simple",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.SimpleResponse("message text").Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "default",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.NewResponse(
					sugo.ResponseDefault,
					"message title",
					"message text",
				).Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "info",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.NewResponse(
					sugo.ResponseInfo,
					"message title",
					"message text",
				).Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "success",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.NewResponse(
					sugo.ResponseSuccess,
					"message title",
					"message text",
				).Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "warning",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.NewResponse(
					sugo.ResponseWarning,
					"message title",
					"message text",
				).Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
		{
			Trigger: "danger",
			Execute: func(req *sugo.Request) (err error) {
				if _, err = req.NewResponse(
					sugo.ResponseDanger,
					"message title",
					"message text",
				).Send(); err != nil {
					return sugo.WrapError(req, err)
				}
				return
			},
		},
	},
}
