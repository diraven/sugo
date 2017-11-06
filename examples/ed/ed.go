package ed

import (
	"github.com/diraven/sugo"
)

// Command contains all ed-related stuff.
var Cmd = &sugo.Command{
	Trigger:            "ed",
	PermittedByDefault: true,
	Description:        "All kinds of Elite: Dangerous related commands.",
	SubCommands: []*sugo.Command{
		System,
	},
}
