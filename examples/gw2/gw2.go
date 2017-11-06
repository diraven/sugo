package gw2

import (
	"github.com/diraven/sugo"
)

// Command contains all ed-related stuff.
var Cmd = &sugo.Command{
	Trigger:            "gw2",
	PermittedByDefault: true,
	Description:        "All kinds of Guild Wars 2 related commands.",
	SubCommands: []*sugo.Command{
		ItemInfo,
	},
}
