package GuildWars2

import (
	"github.com/diraven/sugo"
)

// Module contains all guild wars2-related stuff.
var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "guild_wars2",
		PermittedByDefault: true,
		Description:        "All kinds of Guild Wars 2 related commands.",
		SubCommands: []*sugo.Command{
			itemInfo,
		},
	},
}
