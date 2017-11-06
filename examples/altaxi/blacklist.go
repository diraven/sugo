package altaxi

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"io/ioutil"
	"strconv"
	"strings"
)

type sBlacklistItem struct {
	Author string
	Victim string
}

type sBlacklistStorage struct {
	Items []sBlacklistItem
}

var blacklistStorage = &sBlacklistStorage{}

var TAXI_BLACKLIST_DATA_FILENAME = "taxi_blacklist.json"

func init() {
	blacklistStorage.Items = []sBlacklistItem{}
}

// Command contains all ed-related stuff.
var cmdBlacklist = &sugo.Command{
	Trigger:            "blacklist",
	PermittedByDefault: true,
	Description:        "Allows to manipulate taxi blacklist.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		embed := &discordgo.MessageEmbed{
			Title:  "Blacklisted people",
			Fields: []*discordgo.MessageEmbedField{},
		}

		for i, item := range blacklistStorage.Items {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  strconv.FormatInt(int64(i), 10) + ". " + item.Victim,
				Value: "by " + item.Author,
			})

		}

		_, err = sg.RespondEmbed(m, embed)
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "add",
			Description: "Adds new item into blacklist.",
			Usage:       "@WhoAdded Bad Guy",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				ss := strings.Split(q, " ")
				if len(ss) < 1 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				switch len(ss) {
				case 1: // One-word item name.
					blacklistStorage.Items = append(blacklistStorage.Items, sBlacklistItem{
						m.Author.Username,
						q,
					})
				case 2: // Two-word item name OR one-word item name and mention.
					if len(m.Mentions) > 0 {
						blacklistStorage.Items = append(blacklistStorage.Items, sBlacklistItem{
							m.Mentions[0].Username,
							ss[1],
						})
					} else {
						blacklistStorage.Items = append(blacklistStorage.Items, sBlacklistItem{
							m.Author.Username,
							q,
						})
					}
				case 3: // Two-word item name with "who added" specified as mention.
					if len(m.Mentions) < 1 {
						_, err = sg.RespondBadCommandUsage(m, c, "")
						return
					}
					blacklistStorage.Items = append(blacklistStorage.Items, sBlacklistItem{
						m.Mentions[0].Username,
						strings.TrimSpace(strings.Replace(q, ss[0], "", 1)),
					})
				default:
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}
				_, err = sg.RespondSuccessMention(m, "")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified blacklisted item.",
			Usage:       "1",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				i, err := strconv.ParseInt(q, 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				count := len(blacklistStorage.Items)

				if int(i) >= count || i < 0 {
					_, err = sg.RespondFailMention(m, "Blacklist item \""+q+"\" not found.")
					return
				}

				blacklistStorage.Items = append(blacklistStorage.Items[:i], blacklistStorage.Items[i+1:]...)

				_, err = sg.RespondSuccessMention(m, "")
				if err != nil {
					return
				}
				return
			},
		},
		{
			Trigger:     "check",
			Description: "Check if the person is in the blacklist.",
			Usage:       "Bad Guy",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				embed := &discordgo.MessageEmbed{
					Fields: []*discordgo.MessageEmbedField{},
				}
				var found bool
				for _, item := range blacklistStorage.Items {
					if strings.Index(item.Victim, q) >= 0 {
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:  item.Victim,
							Value: "by " + item.Author,
						})
						found = true
					}
				}
				if found {
					embed.Color = sugo.ColorDanger
					embed.Title = "Blacklisted people found"
					_, err = sg.RespondEmbed(m, embed)
					return
				}

				embed.Color = sugo.ColorSuccess
				embed.Title = "No blacklisted people found"
				_, err = sg.RespondEmbed(m, embed)

				return
			},
		},
	},
	Startup: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Load file.
		data, err := ioutil.ReadFile(TAXI_BLACKLIST_DATA_FILENAME)
		if err != nil {
			return
		}

		// Decode JSON data.
		json.Unmarshal(data, blacklistStorage)
		if err != nil {
			return
		}

		// Log the operation results.

		return
	},
	Teardown: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Encode our data into JSON.
		data, err := json.Marshal(blacklistStorage)
		if err != nil {
			return
		}

		// Save data into file.
		err = ioutil.WriteFile(TAXI_BLACKLIST_DATA_FILENAME, data, 0644)
		if err != nil {
			return
		}

		return
	},
}
