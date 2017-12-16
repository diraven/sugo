package feeds

import (
	"github.com/diraven/sugo"
	"github.com/mmcdole/gofeed"
	"sync"
	"time"
)

var wg sync.WaitGroup

var cPostNewItemsTeardown = make(chan bool)

var fp = gofeed.NewParser()

var feeds = tFeeds{}

func postNewItems(sg *sugo.Instance) {
	wg.Add(1)
	for {
		select {
		case <-cPostNewItemsTeardown: // Bot is shutting down. Exit the loop.
			wg.Done()
			return
		case <-time.After(time.Minute * 15):
			for i, item := range feeds {
				feed, err := fp.ParseURL(item.Url)
				if err == nil {
					for k := len(feed.Items) - 1; k >= 0; k-- {
						if feed.Items[k].PublishedParsed.After(*item.PostedAt) {
							feeds[i].PostedAt = feed.Items[k].PublishedParsed
							sg.ChannelMessageSend(item.ChannelID, "**"+feed.Items[k].Title+"**\n"+feed.Items[k].Link)
						}
					}
					feeds.updatePostedAt(sg, item.ChannelID, item.Url, feeds[i].PostedAt)
				} else {
					sg.ChannelMessageSend(item.ChannelID, item.Url+": "+err.Error())
				}
			}
		}
	}
}

// Module allows to manipulate rss posting settings.
var Module = &sugo.Module{
	Startup:     startup,
	Teardown:    teardown,
	RootCommand: rootCommand,
}
