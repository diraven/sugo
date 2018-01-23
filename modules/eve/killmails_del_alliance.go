package eve

import (
	"github.com/diraven/sugo"
	"strings"
)

var killMailDelAlliance = &sugo.Command{
	Trigger:            "del_alliance",
	PermittedByDefault: true,
	Description:        "Removes alliance from the killmail tracking.",
	Usage:              "01234567890",
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		var err error

		// Make sure there is a query specified.
		if strings.TrimSpace(req.Query) == "" {
			_, err = sg.RespondBadCommandUsage(req, "", "")
			return err
		}

		err = killmails.delAlliance(sg, req.Channel.ID, req.Query)
		if err != nil {
			return err
		}

		_, err = sg.RespondSuccess(req, "", "")
		return err
	},
}
