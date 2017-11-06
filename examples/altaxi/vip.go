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

type sVipItem struct {
	Name string
}

type sVipStorage struct {
	Items []sVipItem
}

var vipStorage = &sVipStorage{}

var TAXI_VIP_DATA_FILENAME = "taxi_vip.json"

func init() {
	vipStorage.Items = []sVipItem{}
}

// Command contains all ed-related stuff.
var cmdVip = &sugo.Command{
	Trigger:            "vip",
	PermittedByDefault: true,
	Description:        "Allows to manipulate taxi VIPs.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		embed := &discordgo.MessageEmbed{
			Title:  "VIPs",
			Fields: []*discordgo.MessageEmbedField{},
		}

		for i, item := range vipStorage.Items {
			embed.Description = embed.Description + strconv.FormatInt(int64(i), 10) + ". " + item.Name + "\n"
		}

		_, err = sg.RespondEmbed(m, embed)
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "add",
			Description: "Adds new VIP.",
			Usage:       "Name Surname",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				ss := strings.Split(q, " ")
				if len(ss) < 1 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}
				vipStorage.Items = append(vipStorage.Items, sVipItem{
					q,
				})
				_, err = sg.RespondSuccessMention(m, "")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified VIP.",
			Usage:       "1",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				i, err := strconv.ParseInt(q, 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				count := len(vipStorage.Items)

				if int(i) >= count || i < 0 {
					_, err = sg.RespondTextMention(m, "VIP \""+q+"\" not found.")
					return
				}

				vipStorage.Items = append(vipStorage.Items[:i], vipStorage.Items[i+1:]...)

				_, err = sg.RespondSuccessMention(m, "")
				if err != nil {
					return
				}
				return
			},
		},
		{
			Trigger:            "check",
			Description:        "Check if the person is VIP.",
			Usage:              "Name Surname",
			PermittedByDefault: true,
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				embed := &discordgo.MessageEmbed{
					Fields: []*discordgo.MessageEmbedField{},
				}
				var found bool
				for _, item := range vipStorage.Items {
					if strings.Index(item.Name, q) >= 0 {
						embed.Description = embed.Description + item.Name + "\n"
						found = true
					}
				}
				if found {
					embed.Color = sugo.ColorSuccess
					embed.Title = "VIPs found"
					_, err = sg.RespondEmbed(m, embed)
					return
				}

				embed.Color = sugo.ColorDanger
				embed.Title = "No VIPs found"
				_, err = sg.RespondEmbed(m, embed)

				return
			},
		},
	},
	Startup: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Load file.
		data, err := ioutil.ReadFile(TAXI_VIP_DATA_FILENAME)
		if err != nil {
			return
		}

		// Decode JSON data.
		json.Unmarshal(data, vipStorage)
		if err != nil {
			return
		}

		return
	},
	Teardown: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Encode our data into JSON.
		data, err := json.Marshal(vipStorage)
		if err != nil {
			return
		}

		// Save data into file.
		err = ioutil.WriteFile(TAXI_VIP_DATA_FILENAME, data, 0644)
		if err != nil {
			return
		}

		return
	},
}
