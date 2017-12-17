package aliases

import (
	"github.com/diraven/sugo"
)

// Help shows help section for appropriate command.
var Module = &sugo.Module{
	Startup:               startup,
	RootCommand:           rootCommand,
	OnBeforeCommandSearch: onBeforeCommandSearch,
}
