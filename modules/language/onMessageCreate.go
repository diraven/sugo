package language

import (
	"github.com/diraven/sugo"
	"github.com/nicksnyder/go-i18n/i18n"
)

var translateFuncCache = map[string]*i18n.TranslateFunc{}

func onMessageCreate(sg *sugo.Instance, req *sugo.Request) error {
	// Get current user language.
	language, err := languages.get(sg, req)
	if err != nil {
		return err
	}

	// Try to get translate func from cache.
	_, ok := translateFuncCache[language]
	if !ok {
		// Create new translation function for given language.
		tfunc, err := i18n.Tfunc(language, "en")
		if err != nil {
			return err
		}
		// Cache newly created translation function.
		translateFuncCache[language] = &tfunc
	}

	// Set request translation function depending on user language.
	req.TranslateFunc = translateFuncCache[language]

	return nil
}
