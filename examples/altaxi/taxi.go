package altaxi

import (
	"github.com/diraven/sugo"
)

// Command contains all alttaxi-related stuff.
var Cmd = &sugo.Command{
	Trigger:            "taxi",
	PermittedByDefault: true,
	Description:        "Contains taxi-related  commands..",
	SubCommands: []*sugo.Command{
		cmdOrder,
		cmdReport,
		cmdBlacklist,
		cmdLocations,
		cmdVip,
	},
}
