package stats

import (
	"github.com/diraven/sugo"
)

var stats = tStats{}

// Module allows to manipulate rss posting settings.
var Module = &sugo.Module{
	Startup:     startup,
	RootCommand: rootCommand,
	OnPresenceUpdate: onPresenceUpdate,
}
