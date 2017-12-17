package elite_dangerous

import (
	"github.com/diraven/sugo"
)

// Command contains all ed-related stuff.

var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:            "elite_dangerous",
		PermittedByDefault: true,
		Description:        "All kinds of Elite: Dangerous related commands.",
		SubCommands: []*sugo.Command{
			factions,
		},
	},
}
