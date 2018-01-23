package eve

import (
	"github.com/diraven/sugo"
)

var killMail = &sugo.Command{
	Trigger:     "killmails",
	RootOnly:    true,
	Description: "Provides killmail-related facilities.",
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		// Get all alliances.
		//resp := "Alliances: "
		//for allianceID, channelID := range killmails.allianceIDs {
		//	if req.Channel.ID == channelID {
		//		resp = resp + " " + allianceID
		//	}
		//}

		// Get all corporations.
		resp := "Corporations: "
		for corporationID, channelID := range killmails.corporationIDs {
			if req.Channel.ID == channelID {
				resp = resp + " " + corporationID
			}
		}

		// Return response.
		_, err := sg.RespondInfo(req, "Subscribed killmail IDs", resp)

		return err
	},
	SubCommands: []*sugo.Command{
		//killMailAddAlliance,
		killMailAddCorporation,
		//killMailDelAlliance,
		killMailDelCorporation,
	},
}
