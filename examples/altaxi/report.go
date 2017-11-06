package altaxi

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"io/ioutil"
	"strconv"
	"time"
)

type reportItem struct {
	Count  int
	Amount int
}

type sReportStorage struct {
	Items map[string]reportItem
}

var reportStorage = &sReportStorage{}
var loc *time.Location
var reportDateTime time.Time

var TAXI_REPORT_DATA_FILENAME = "taxi_report.json"

func init() {
	var err error
	loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	reportDateTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	reportStorage.Items = make(map[string]reportItem)
}

func postReport(x context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
	if len(reportStorage.Items) > 0 {
		embed := &discordgo.MessageEmbed{
			Title:  "Report",
			Fields: []*discordgo.MessageEmbedField{},
		}

		var authorUser *discordgo.Member
		var guild *discordgo.Guild
		guild, err = sg.GuildFromMessage(m)
		if err != nil {
			return
		}
		for authorID, reportItem := range reportStorage.Items {
			authorUser, err = sg.State.Member(guild.ID, authorID)
			if err != nil {
				return
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: authorUser.User.Username,
				Value: strconv.FormatInt(int64(reportItem.Amount), 10) +
					" (" +
					strconv.FormatInt(int64(reportItem.Count), 10) +
					")",
				Inline: true,
			})
		}
		_, err = sg.RespondEmbed(m, embed)
	}

	return
}

var cmdReport = &sugo.Command{
	Trigger:            "report",
	PermittedByDefault: true,
	Description:        "Allows to report and calculate daily earnings..",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		err = postReport(ctx, c, q, sg, m)
		return
	},
	SubCommands: []*sugo.Command{
		cmdReportAdd,
		cmdReportReset,
	},
	Startup: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Load file.
		data, err := ioutil.ReadFile(TAXI_REPORT_DATA_FILENAME)
		if err != nil {
			return
		}

		// Decode JSON data.
		json.Unmarshal(data, reportStorage)
		if err != nil {
			return
		}

		return
	},
	Teardown: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Encode our data into JSON.
		data, err := json.Marshal(reportStorage)
		if err != nil {
			return
		}

		// Save data into file.
		err = ioutil.WriteFile(TAXI_REPORT_DATA_FILENAME, data, 0644)
		if err != nil {
			return
		}

		return
	},
}
