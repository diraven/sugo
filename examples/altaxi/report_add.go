package altaxi

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"strconv"
	"strings"
	"time"
)

var cmdReportAdd = &sugo.Command{
	Trigger:            "add",
	PermittedByDefault: true,
	Description:        "Adds earnings to the daily report.",
	Usage:              "50k",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {

		duration := time.Since(reportDateTime)

		if duration.Hours() >= 24 {
			err = postReport(ctx, c, q, sg, m)
			if err != nil {
				return
			}
			_, err = sg.RespondTextMention(m, "New day, new report. All yesterday's data discarded!")
			if err != nil {
				return
			}
			now := time.Now().In(loc)
			reportDateTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			reportStorage.Items = make(map[string]reportItem)
		}

		//if q == "" {
		//	order, exists := pendingOrders[m.Author.ID]
		//	if exists {
		//		q = strconv.FormatInt(order.amount, 10)
		//		delete(pendingOrders, m.Author.ID)
		//	} else {
		//		sg.RespondTextMention(m, "You have no pending orders.")
		//		return
		//	}
		//}

		ksCount := strings.Count(q, "k")

		q = strings.Replace(q, "k", "", -1)

		amount, err := strconv.ParseFloat(q, 64)
		if err != nil {
			_, err = sg.RespondBadCommandUsage(m, c, "")
			return
		}

		for i := 0; i < ksCount; i++ {
			amount = amount * 1000
		}

		// Store Amount.
		item, exists := reportStorage.Items[m.Author.ID]
		if exists {
			item.Count++
			item.Amount = item.Amount + int(amount)
			reportStorage.Items[m.Author.ID] = item
		} else {
			reportStorage.Items[m.Author.ID] = reportItem{
				Count:  1,
				Amount: int(amount),
			}
		}

		// Respond to user.
		_, err = sg.RespondSuccessMention(m, "")
		return
	},
}
