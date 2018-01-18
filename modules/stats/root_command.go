package stats

import (
	"github.com/diraven/sugo"
	"strconv"
)

var rootCommand = &sugo.Command{
	Trigger:            "stats",
	Description:        "Gives general server stats.",
	PermittedByDefault: true,
	Execute: func(sg *sugo.Instance, req *sugo.Request) error {
		if _, err := sg.RespondNotImplemented(req); err != nil {
			return err
		}

		return nil
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:            "playing",
			Description:        "Shows stats about games played most.",
			PermittedByDefault: true,
			//Usage:         "http://example.com/rss/",
			//AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var response string

				gamesNames, err := stats.getMostPlayedGames(sg, req.Guild.ID)
				if err != nil {
					return err
				}

				for i, gameName := range gamesNames {
					response = response + strconv.Itoa(i+1) + ". " + gameName + "\n"
				}

				if _, err := sg.RespondInfo(req, "most played games", response); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Trigger:            "messaging",
			Description:        "Shows stats about most active users.",
			PermittedByDefault: true,
			//Usage:         "http://example.com/rss/",
			//AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var response string

				users, err := stats.getMostMessagingUsers(sg, req.Guild.ID)
				if err != nil {
					return err
				}

				for i, user := range users {
					response = response + strconv.Itoa(i+1) + ". " + user.Mention() + "\n"
				}

				if _, err := sg.RespondInfo(req, "most active users", response); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
