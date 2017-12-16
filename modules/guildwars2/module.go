package guildwars2

import (
	"github.com/diraven/sugo"
)

// Module contains all guildwars2-related stuff.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "guildwars2",
		PermittedByDefault: true,
		Description:        "All kinds of Guild Wars 2 related commands.",
		SubCommands: []*sugo.Command{
			itemInfo,
		},
	},
}