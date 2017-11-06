package altaxi

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
	"math"
	"strconv"
	"strings"
)

type order struct {
	distance int64
}

type price struct {
	name  string
	perKM int64
}

var pendingOrders = make(map[string]order)

// Command contains all ed-related stuff.
var cmdOrder = &sugo.Command{
	Trigger:            "order",
	PermittedByDefault: true,
	Description:        "Calculates travel distance and price.",
	Usage:              "123123 123123 Name Surname",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		ss := strings.Split(q, " ")
		if len(ss) < 2 {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		// Replace locations if any found.
		var replaced bool
		for _, v := range locationsStorage.Items {
			if ss[0] == v.Name || ss[1] == v.Name {
				q = strings.Replace(q, v.Name, v.Coordinates, 1)
				replaced = true
			}
		}

		if replaced {
			ss = strings.Split(q, " ")
		}

		startCoords := ss[0]
		q = helpers.ConsumePrefix(q, startCoords)
		if len(startCoords) < 6 {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		finishCoords := ss[1]
		q = helpers.ConsumePrefix(q, finishCoords)
		if len(finishCoords) < 6 {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		x1, err := strconv.ParseInt(startCoords[:3], 10, 0)
		if err != nil {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}
		y1, err := strconv.ParseInt(startCoords[3:], 10, 0)
		if err != nil {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}
		x2, err := strconv.ParseInt(finishCoords[:3], 10, 0)
		if err != nil {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}
		y2, err := strconv.ParseInt(finishCoords[3:], 10, 0)
		if err != nil {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		distance := math.Sqrt(
			math.Pow(float64(x2)-float64(x1), 2)+
				math.Pow(float64(y2)-float64(y1), 2),
		) / 10 // in 100m s

		prices := []price{
			{"Land", 7000},
			{"Air", 8000},
		}

		embed := &discordgo.MessageEmbed{
			Title: "Order",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Distance",
					Value: strconv.FormatFloat(distance, 'f', 2, 64) + "km",
				},
			},
		}

		for _, v := range prices {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   v.name,
				Value:  strconv.FormatInt(int64(math.Ceil(distance*float64(v.perKM)/1000)), 10) + "k",
				Inline: true,
			})
		}

		_, err = sg.RespondEmbed(m, embed)
		if err != nil {
			return
		}

		pendingOrders[m.Author.ID] = order{distance: int64(distance)}
		return
	},
	SubCommands: []*sugo.Command{
		orderCancel,
	},
}
