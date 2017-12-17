package triggers

import (
	"github.com/diraven/sugo"
)

var triggers = tTriggers{}

// Module allows to set custom bot trigger.
var Module = &sugo.Module{
	Startup:            startup,
	RootCommand:        rootCommand,
	OnBeforeBotTriggerDetect: onBeforeBotTriggerDetect,
}
