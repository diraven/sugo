package eve

import (
	"github.com/diraven/sugo"
	"sync"
	"time"
	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"strconv"
)

var wg sync.WaitGroup
var cPostNewItemsTeardown = make(chan bool)
var killmails = tKillmails{}

func postKillmails(sg *sugo.Instance) {
	wg.Add(1)
	for {
		select {
		case <-cPostNewItemsTeardown: // Bot is shutting down. Exit the loop.
			wg.Done()
			return
		case <-time.After(time.Millisecond):
			resp, err := http.Get("https://redisq.zkillboard.com/listen.php")
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println(fmt.Errorf("unexpected http GET status: %s", resp.Status))
			}
			var killmail = Killmail{}
			err = json.NewDecoder(resp.Body).Decode(&killmail)
			if err != nil {
				log.Println(fmt.Errorf("cannot decode JSON: %v", err))
			}
			channelID, ok := killmails.corporationIDs[strconv.Itoa(killmail.Package.Killmail.Victim.CorporationID)]
			if ok {
				sg.ChannelMessageSend(channelID, "https://zkillboard.com/kill/"+strconv.Itoa(killmail.Package.KillID)+"/")
			} else {
				for _, attacker := range killmail.Package.Killmail.Attackers {
					channelID, ok := killmails.corporationIDs[strconv.Itoa(attacker.CorporationID)]
					if ok {
						sg.ChannelMessageSend(channelID, "https://zkillboard.com/kill/"+strconv.Itoa(killmail.Package.KillID)+"/")
					}
				}
			}
		}
	}
}

// Module contains all eve-related stuff.
var Module = &sugo.Module{
	Startup: startup,
	RootCommand: &sugo.Command{
		Trigger:            "eve",
		PermittedByDefault: true,
		Description:        "All kinds of EVE Online related commands.",
		SubCommands: []*sugo.Command{
			killMail,
		},
	},
}
