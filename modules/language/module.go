package language

import (
	"github.com/diraven/sugo"
)

const defaultLanguage = "en"

// AvailableLanguages contains selection of all languages users can select via respective command.
var AvailableLanguages = []string{
	"en",
}

var languages = tLanguages{}

// Module allows for language setting.
var Module = &sugo.Module{
	RootCommand:     rootCommand,
	Startup:         startup,
	OnMessageCreate: onMessageCreate,
}
