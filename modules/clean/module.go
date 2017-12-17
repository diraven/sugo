package clean

import (
	"github.com/diraven/sugo"
	"time"
	"github.com/bwmarrin/discordgo"
	"strings"
	"strconv"
	"github.com/diraven/sugo/helpers"
	"context"
)

var Module = &sugo.Module{
	RootCommand: &sugo.Command{
		Trigger:             "clean",
		AllowDefaultChannel: true,
		ParamsAllowed:       true,
		Description:         "Removes last n messages, or last n messages of the given @user (if specified) (100 max).",
		Usage:               "[@user] [messages_count]",
		Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
			var err error

			// Command has to have 1 or 2 parameters.
			ss := strings.Split(q, " ")

			var batchSize = 100  // Amount of messages to get in one go.
			var maxCount = 100   // Maximum amount of messages deleted.
			var userID string    // UserID to delete messages for.
			var count = maxCount // Default amount of messages to be deleted.

			switch len(ss) {
			case 1: // Means we have got either user mention or amount of messages to delete.
				if ss[0] == "" { // No parameters given.
					break
				}

				if len(m.Mentions) > 0 { // Get user mention if available.
					// Get user id.
					userID = m.Mentions[0].ID
				} else { // Get amount of messages to delete.
					count, err = strconv.Atoi(ss[0]) // Try to parse count.
					if err != nil {
						return err
					}
				}
				break
			case 2: // Means we've got both user mention and amount of messages to delete.
				if len(m.Mentions) == 0 { // Query must have mention.
					_, err = sg.RespondBadCommandUsage(m, c, "")
					if err != nil {
						return err
					}
				}
				userID = m.Mentions[0].ID

				// Try to get count of messages to delete.
				count, err = strconv.Atoi(ss[0]) // Try first argument.
				if err != nil { // If first argument did not work.
					count, err = strconv.Atoi(ss[1]) // Try second one.
					if err != nil {
						_, err = sg.RespondBadCommandUsage(m, c, "")
						if err != nil {
							return err
						}
					}
				}
				break
			default:
				_, err = sg.RespondBadCommandUsage(m, c, "")
				if err != nil {
					return err
				}
				return err
			}

			// Validate count.
			if count > maxCount {
				_, err = sg.RespondBadCommandUsage(m, c, "Max messages count is "+strconv.Itoa(maxCount)+".")
				if err != nil {
					return err
				}
			}

			lastMessageID := m.ID                // To store last message id.
			var tmpMessages []*discordgo.Message // To store 100 current messages that are being scanned.
			var messageIDs []string              // Resulting slice of messages to  be deleted.
			limit := batchSize                   // Default limit per batch.

			if userID == "" && count < batchSize { // If user ID is not specified - we retrieve and delete exact count of messages specified.
				limit = count
			}

			// Start getting messages.
		messageLoop:
			for {
				// Get next 100 messages.
				tmpMessages, err = sg.ChannelMessages(m.ChannelID, limit, lastMessageID, "", "")
				if err != nil {
					return err
					break messageLoop
				}

				// For each message.
				for _, message := range tmpMessages {
					// Get message creation date.
					var then time.Time
					then, err = helpers.DiscordTimestampToTime(string(message.Timestamp))
					if err != nil {
						return err
					}

					if time.Since(then).Hours() >= 24*14 {
						// We are unable to delete messages older then 14 days.
						_, err = sg.RespondDanger(m, "unable to delete messages older then 2 weeks")
						if err != nil {
							return err
						}
						break messageLoop
					}
					if userID != "" {
						// If user ID is specified, we compare message with the user ID.
						if message.Author.ID == userID {
							messageIDs = append(messageIDs, message.ID)
						}
					} else {
						// Otherwise just add message ID to the list for deletion.
						messageIDs = append(messageIDs, message.ID)
					}

					// If we have enough messages staged for deletion.
					if len(messageIDs) >= count {
						// Finish looking for messages.
						break messageLoop
					}

				}

				if len(tmpMessages) < batchSize {
					break messageLoop // We have no messages left to scan.
				}

				if len(tmpMessages) == batchSize {
					lastMessageID = tmpMessages[batchSize-1].ID // Next time start scanning from the message specified.
				}
			}

			// Delete command itself. Ignore errors (such as message already deleted by someone) for now.
			_ = sg.ChannelMessageDelete(m.ChannelID, m.ID)

			// Perform selected messages deletion. Ignore errors (such as message already deleted by someone) for now.
			_ = sg.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)

			// Notify user about deletion.
			msg, err := sg.RespondSuccess(m, "done, this message will self-destruct in 10 seconds")
			if err != nil {
				return err
			}

			// Wait for 10 seconds.
			time.Sleep(10 * time.Second)

			// Delete notification. Ignore errors (such as message already deleted by someone) for now.
			_ = sg.ChannelMessageDelete(msg.ChannelID, msg.ID)
			return err
		},
	},
}