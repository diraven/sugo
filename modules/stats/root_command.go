package stats

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strconv"
)

var rootCommand = &sugo.Command{
	Trigger:            "stats",
	Description:        "Gives general server stats.",
	PermittedByDefault: true,
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		if _, err := sg.RespondNotImplemented(m); err != nil {
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
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var response string

				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				gamesNames, err := stats.getMostPlayedGames(sg, guild.ID)
				if err != nil {
					return err
				}

				for i, gameName := range gamesNames {
					response = response + strconv.Itoa(i+1) + ". " + gameName + "\n"
				}

				if _, err := sg.RespondInfo(m, "most played games", response); err != nil {
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
			Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var response string

				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				users, err := stats.getMostMessagingUsers(sg, guild.ID)
				if err != nil {
					return err
				}

				for i, user := range users {
					response = response + strconv.Itoa(i+1) + ". " + user.Mention() + "\n"
				}

				if _, err := sg.RespondInfo(m, "most active users", response); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
