package language

import (
	"github.com/diraven/sugo"
	"strings"
)

var rootCommand = &sugo.Command{
	Trigger:             "language",
	PermittedByDefault:  true,
	AllowDefaultChannel: true,
	AllowParams:         true,
	Description:         "Shows or sets language.",
	//Usage:               defaultFormat,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Check query.
		if req.Query == "" {
			// Query is empty, get current user language
			language, err := languages.get(sg, req)
			if err != nil {
				return err
			}

			// Respond with the resulting language to user.
			if _, err := sg.RespondInfo(req, "your current language", language); err != nil {
				return err
			}

			// We are done here.
			return nil
		}

		if req.Query == "reset" {
			// Settings reset requested.
			if err := languages.reset(sg, req.Message.Author.ID); err != nil {
				return err
			}

			// Notify user about successfull operation.
			if _, err := sg.RespondSuccess(req, "", ""); err != nil {
				return err
			}

			// We are done here.
			return nil
		}

		// Query is not empty, try to find language among available languages.
		for _, v := range AvailableLanguages {
			if v == req.Query {
				// We have found that language requested by user is available.
				// Save user setting.
				languages.set(sg, req.Message.Author.ID, req.Query)

				// Notify user about successful operation.
				_, err := sg.RespondSuccess(req, "", "")
				if err != nil {
					return err
				}

				// We are done here.
				return nil
			}
		}

		// We did not find requested language in supported languages list.
		// Notify user about failure.
		_, err := sg.RespondBadCommandUsage(req, "", req.Query+" is not supported, supported languages are: "+strings.Join(AvailableLanguages, ","))
		if err != nil {
			return err
		}

		// We are done here.
		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:             "guild",
			AllowDefaultChannel: true,
			AllowParams:         true,
			RootOnly:            true,
			Description:         "Shows or sets guild timezone.",
			Usage:               "en",
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				// Check query.
				if req.Query == "" {
					if _, err := sg.RespondBadCommandUsage(req, "", ""); err != nil {
						return err
					}
					// We are done here.
					return nil
				}

				if req.Query == "reset" {
					// Settings reset requested.
					if err := languages.reset(sg, req.Guild.ID); err != nil {
						return err
					}

					// Notify user about successfull operation.
					if _, err := sg.RespondSuccess(req, "", ""); err != nil {
						return err
					}

					// We are done here.
					return nil
				}

				// Query is not empty, try to find language among available languages.
				for _, v := range AvailableLanguages {
					if v == req.Query {
						// We have found that language requested by user is available.
						// Save user setting.
						languages.set(sg, req.Guild.ID, req.Query)

						// Notify user about successful operation.
						_, err := sg.RespondSuccess(req, "", "")
						if err != nil {
							return err
						}

						// We are done here.
						return nil
					}
				}

				// We did not find requested language in supported languages list.
				// Notify user about failure.
				if _, err := sg.RespondBadCommandUsage(req, "", req.Query+" is not supported, supported languages are: "+strings.Join(AvailableLanguages, ",")); err != nil {
					return err
				}

				// We are done here.
				return nil
			},
		},
	},
}
