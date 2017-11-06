package std

import (
	"github.com/diraven/sugo"
)

// Test is just a testing command
var Test = &sugo.Command{
	Trigger:             "test",
	PermittedByDefault:  true,
	Description:         "Test command.",
	TextResponse:        "Test passed.",
	AllowDefaultChannel: true,
	SubCommands: []*sugo.Command{
		{
			Trigger:            "subtest1",
			PermittedByDefault: true,
			Description:        "subTest1 command.",
			TextResponse:       "subTest1 passed.",
			SubCommands: []*sugo.Command{
				{
					Trigger:            "subtest11",
					PermittedByDefault: true,
					Description:        "subTest11 command.",
					TextResponse:       "subTest11 passed.",
				},
				{
					Trigger:            "subtest12",
					PermittedByDefault: true,
					Description:        "subTest12 command.",
					TextResponse:       "subTest12 passed.",
				},
			},
		},
		{
			Trigger:            "subtest2",
			PermittedByDefault: true,
			Description:        "subTest2 command.",
			TextResponse:       "subTest2 passed.",
		},
	},
}
