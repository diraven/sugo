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

type sLocation struct {
	Name        string
	Coordinates string
}

type sLocationsStorage struct {
	Items []sLocation
}

var locationsStorage = &sLocationsStorage{}

var TAXI_LOCATIONS_DATA_FILENAME = "taxi_locations.json"

func init() {
	locationsStorage.Items = []sLocation{}
}

// Command contains all ed-related stuff.
var cmdLocations = &sugo.Command{
	Trigger:            "locations",
	PermittedByDefault: true,
	Description:        "Allows to manipulate taxi locations.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		embed := &discordgo.MessageEmbed{
			Title: "Locations",
		}

		for i, item := range locationsStorage.Items {
			embed.Description = embed.Description +
				strconv.FormatInt(int64(i), 10) +
				". " +
				item.Name +
				": " +
				item.Coordinates +
				"\n"
		}

		_, err = sg.RespondEmbed(m, embed)
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "add",
			Description: "Adds new location.",
			Usage:       "locationname 123123",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				ss := strings.Split(q, " ")
				if len(ss) < 2 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				_, err = strconv.ParseInt(ss[1], 10, 64)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				if len(ss[1]) != 6 {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				locationsStorage.Items = append(locationsStorage.Items, sLocation{
					ss[0],
					ss[1],
				})

				_, err = sg.RespondSuccessMention(m, "")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified shortcut.",
			Usage:       "1",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				i, err := strconv.ParseInt(q, 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				count := len(locationsStorage.Items)

				if int(i) >= count || i < 0 {
					_, err = sg.RespondTextMention(m, "Locations item \""+q+"\" not found.")
					return
				}

				locationsStorage.Items = append(locationsStorage.Items[:i], locationsStorage.Items[i+1:]...)

				_, err = sg.RespondSuccessMention(m, "")
				if err != nil {
					return
				}
				return
			},
		},
		{
			Trigger:     "get",
			Description: "Get location.",
			Usage:       "kavala",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				if q == "" {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				for _, v := range locationsStorage.Items {
					if strings.Index(v.Name, q) > -1 {
						_, err = sg.RespondTextMention(m, v.Coordinates)
						return
					}
				}
				_, err = sg.RespondTextMention(m, "Location not found...")
				return
			},
		},
	},
	Startup: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Load file.
		data, err := ioutil.ReadFile(TAXI_LOCATIONS_DATA_FILENAME)
		if err != nil {
			return
		}

		// Decode JSON data.
		json.Unmarshal(data, locationsStorage)
		if err != nil {
			return
		}

		return
	},
	Teardown: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Encode our data into JSON.
		data, err := json.Marshal(locationsStorage)
		if err != nil {
			return
		}

		// Save data into file.
		err = ioutil.WriteFile(TAXI_LOCATIONS_DATA_FILENAME, data, 0644)
		if err != nil {
			return
		}

		return
	},
}
