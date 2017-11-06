package ed

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type timeoutError struct {
	s string
}

func (e timeoutError) Error() string {
	return "Server took too long to respond. Operation aborted."
}

// System returns info about given system.
var System = &sugo.Command{
	Timeout:            5 * time.Second,
	Trigger:            "factions",
	PermittedByDefault: true,
	Description:        "Provides minor factions info on the given system.",
	Usage:              "Solar System Name",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Make sure there is a query specified.
		if strings.TrimSpace(q) == "" {
			_, err = sg.RespondTextMention(m, "Please, specify the system to show data for. See \""+c.FullHelpPath(sg)+"\" for details.")
			return
		}

		// Get system ID by searching with query q.
		var systemID int
		systemID, err = getSystemID(ctx, q)
		if err != nil {
			if _, ok := err.(timeoutError); ok {
				_, err = sg.RespondTextMention(m, fmt.Sprint(err))
				if err != nil {
					return
				}
			}
			return
		}
		// If there are no systems found - we will end up with systemID=0
		if systemID == 0 {
			_, err = sg.RespondTextMention(m, "Oops... Looks like no systems found.")
			if err != nil {
				return
			}
			return
		}

		// Build embed
		var embed *discordgo.MessageEmbed
		embed, err = getSystemEmbed(ctx, systemID)
		if err != nil {
			return
		}
		// It is possible that embed is nil if ctx timeout is reached.
		if embed == nil {
			if _, ok := err.(timeoutError); ok {
				_, err = sg.RespondTextMention(m, fmt.Sprint(err))
				if err != nil {
					return
				}
			}
			return
		}

		_, err = sg.RespondEmbed(m, embed)
		if err != nil {
			return
		}
		return
	},
}

type sSystem struct {
	ID    int    `json:"id"`
	Title string `json:"name"`
}

var searchURLFormat = "https://eddb.io/system/search?system[name]=%s"
var systemURLFormat = "https://eddb.io/system/factions/%d"

func getSystemID(ctx context.Context, q string) (systemID int, err error) {
	// Generate search url for the given query string.
	urlString := fmt.Sprintf(searchURLFormat, url.QueryEscape(q))
	log.Println(urlString)

	// Prepare new request.
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return
	}

	// Add context to the request (it will take care of timeout if one set).
	r, err := http.DefaultClient.Do(req.WithContext(ctx))

	// It is possible that timeout is already reached by the moment context is added into http client.
	// In that case request will be nil. We have to check for this.
	if r != nil {
		defer r.Body.Close()
	} else {
		err = timeoutError{}
		return
	}
	if err != nil {
		return
	}

	// Decode response.
	var systems []sSystem
	err = json.NewDecoder(r.Body).Decode(&systems)
	if err != nil {
		return
	}

	// Check how many systems we have got.
	if len(systems) > 0 {
		// At least one system is found.
		systemID = systems[0].ID
	}

	return
}

func getSystemEmbed(ctx context.Context, id int) (embed *discordgo.MessageEmbed, err error) {
	// Generate url to get system data from.
	urlString := fmt.Sprintf(systemURLFormat, id)

	// Build request with the url generated.
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return
	}

	// Supply request with context.
	r, err := http.DefaultClient.Do(req.WithContext(ctx))
	// It is possible that timeout is reached by this moment, so we need to check if r is not nil to avoid panic.
	if r != nil {
		defer r.Body.Close()
	} else {
		err = timeoutError{}
		return
	}
	if err != nil {
		return
	}

	// Parse the server response.
	d, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return
	}

	// Build system embed.
	embed = &discordgo.MessageEmbed{
		URL:    urlString,
		Title:  d.Find("title").Text(),
		Fields: []*discordgo.MessageEmbedField{},
	}

	// Fill in influence fields.
	d.Find(".systemFactionRow").Each(func(k int, v *goquery.Selection) {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   strings.TrimSpace(v.Find(".factionName").Text()),
			Value:  strings.TrimSpace(v.Find(".factionInfluence").Text()),
			Inline: true,
		})
	})

	return
}
