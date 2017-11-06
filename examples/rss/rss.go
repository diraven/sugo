package rss

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type sItem struct {
	ChannelID  string
	Url        string
	LastUpdate *time.Time
}

type sStorage struct {
	Items []sItem
}

var storage = &sStorage{}

var DATA_FILENAME = "rss.json"

func init() {
	storage.Items = []sItem{}
}

var wg sync.WaitGroup

var cPostNewItemsTeardown = make(chan bool)

var fp = gofeed.NewParser()

func postNewItems(sg *sugo.Instance) {
	wg.Add(1)
	for {
		select {
		case <-cPostNewItemsTeardown: // Bot is shutting down. Exit the loop.
			wg.Done()
			return
		case <-time.After(time.Minute * 15):
			for i, item := range storage.Items {
				feed, err := fp.ParseURL(item.Url)
				if err == nil {
					for k := len(feed.Items) - 1; k >= 0; k-- {
						if feed.Items[k].PublishedParsed.After(*item.LastUpdate) {
							storage.Items[i].LastUpdate = feed.Items[k].PublishedParsed
							sg.ChannelMessageSend(item.ChannelID, "**"+feed.Items[k].Title+"**\n"+feed.Items[k].Link)
						}
					}
				} else {
					sg.ChannelMessageSend(item.ChannelID, strconv.Itoa(i) + ". " + err.Error())
				}
			}
		}
	}
}

// CmdRSS allows to manipulate rss posting settings.
var Cmd = &sugo.Command{
	Trigger:     "rss",
	Description: "Allows to manipulate RSS postings.",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		embed := &discordgo.MessageEmbed{
			Title:  "RSSs",
			Fields: []*discordgo.MessageEmbedField{},
		}

		for i, item := range storage.Items {
			embed.Description = embed.Description + strconv.FormatInt(int64(i), 10) + ". " + item.Url + "\n"
		}

		_, err = sg.RespondEmbed(m, embed)
		return
	},
	SubCommands: []*sugo.Command{
		{
			Trigger:     "add",
			Description: "Adds new RSS.",
			Usage:       "http://example.com/rss/",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				// Validate RSS url.
				_, err = fp.ParseURL(q)
				if err != nil {
					_, err = sg.RespondFailMention(m, err.Error())
					return
				}

				// Add requested URL to the list.
				now := time.Now()
				storage.Items = append(storage.Items, sItem{
					m.ChannelID,
					q,
					&now,
				})

				_, err = sg.RespondSuccessMention(m, "")
				return
			},
		},
		{
			Trigger:     "del",
			Description: "Deletes specified RSS.",
			Usage:       "1",
			Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
				i, err := strconv.ParseInt(q, 10, 0)
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					return
				}

				count := len(storage.Items)

				if int(i) >= count || i < 0 {
					_, err = sg.RespondFailMention(m, "RSS \""+q+"\" not found.")
					return
				}

				storage.Items = append(storage.Items[:i], storage.Items[i+1:]...)

				_, err = sg.RespondSuccessMention(m, "")
				if err != nil {
					return
				}
				return
			},
		},
	},
	Startup: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		// Check if file exists.
		if _, err = os.Stat(DATA_FILENAME); err == nil {
			// Load file.
			data, err := ioutil.ReadFile(DATA_FILENAME)
			if err != nil {
				return err
			}

			// Decode JSON data.
			json.Unmarshal(data, storage)
			if err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			// If there are any errors other then "file does not exist" - report error and shutdown.
			return
		}

		go postNewItems(sg)

		return nil

	},
	Teardown: func(c *sugo.Command, sg *sugo.Instance) (err error) {
		cPostNewItemsTeardown <- true
		wg.Wait()

		// Encode our data into JSON.
		data, err := json.Marshal(storage)
		if err != nil {
			return
		}

		// Save data into file.
		err = ioutil.WriteFile(DATA_FILENAME, data, 0644)
		if err != nil {
			return
		}

		return
	},
}
