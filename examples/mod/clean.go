package mod

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/helpers"
	"strconv"
	"strings"
	"time"
)

// Clean removes last n messages, or last n messages of the given @user (if specified) (100 max).
var Clean = &sugo.Command{
	Trigger: "clean",
	//RootOnly: true,
	//PermittedByDefault: true,
	AllowDefaultChannel: true,
	Description:         "Removes last n messages, or last n messages of the given @user (if specified) (100 max).",
	Usage:               "[@user] [messages_count]",
	Execute: func(ctx context.Context, c *sugo.Command, q string, sg *sugo.Instance, m *discordgo.Message) (err error) {
		// Command has to have 1 or 2 parameters.
		ss := strings.Split(q, " ")

		var batchSize = 100      // Amount of messages to get in one go.
		var maxCount = 100       // Maximum amount of messages deleted.
		var userID string        // UserID to delete messages for.
		var count int = maxCount // Default amount of messages to be deleted.

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
					return
				}
			}
			break
		case 2: // Means we've got both user mention and amount of messages to delete.
			if len(m.Mentions) == 0 { // Query must have mention.
				_, err = sg.RespondBadCommandUsage(m, c, "")
				if err != nil {
					return
				}
			}
			userID = m.Mentions[0].ID

			// Try to get count of messages to delete.
			count, err = strconv.Atoi(ss[0]) // Try first argument.
			if err != nil {                  // If first argument did not work.
				count, err = strconv.Atoi(ss[1]) // Try second one.
				if err != nil {
					_, err = sg.RespondBadCommandUsage(m, c, "")
					if err != nil {
						return
					}
				}
			}
			break
		default:
			_, err = sg.RespondBadCommandUsage(m, c, "")
			if err != nil {
				return
			}
			return
		}

		// Validate count.
		if count > maxCount {
			_, err = sg.RespondBadCommandUsage(m, c, "Max messages count is "+strconv.Itoa(maxCount)+".")
			if err != nil {
				return
			}
		}

		last_message_id := m.ID                // To store last message id.
		tmp_messages := []*discordgo.Message{} // To store 100 current messages that are being scanned.
		messageIDs := []string{}               // Resulting slice of messages to  be deleted.
		limit := batchSize                     // Default limit per batch.

		if userID == "" && count < batchSize { // If user ID is not specified - we retreive and delete exact count of messages specified.
			limit = count
		}

		// Start getting messages.
	message_loop:
		for {
			// Get next 100 messages.
			tmp_messages, err = sg.ChannelMessages(m.ChannelID, limit, last_message_id, "", "")
			if err != nil {
				return
				break message_loop
			}

			// For each message.
			for _, message := range tmp_messages {
				// Get message creation date.
				var then time.Time
				then, err = helpers.DiscordTimestampToTime(string(message.Timestamp))
				if err != nil {
					return
				}

				if time.Since(then).Hours() >= 24*14 {
					// We are unable to delete messages older then 14 days.
					break message_loop
					_, err = sg.RespondFailMention(m, "Unfortunately I'm unable to delete messages older then 2 weeks.")
					if err != nil {
						return
					}
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
					break message_loop
				}

			}

			if len(tmp_messages) < batchSize {
				break message_loop // We have no messages left to scan.
			}

			if len(tmp_messages) == batchSize {
				last_message_id = tmp_messages[batchSize-1].ID // Next time start scanning from the message specified.
			}
		}

		// Delete command itself. Ignore errors (such as message already deleted by someone) for now.
		_ = sg.ChannelMessageDelete(m.ChannelID, m.ID)

		// Perform selected messages deletion. Ignore errors (such as message already deleted by someone) for now.
		_ = sg.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)

		// Notify user about deletion.
		mymsg, err := sg.RespondSuccessMention(m, "Done. This message will self-destruct in 10 seconds.")
		if err != nil {
			return
		}

		// Wait for 10 seconds.
		time.Sleep(10 * time.Second)

		// Delete notification. Ignore errors (such as message already deleted by someone) for now.
		_ = sg.ChannelMessageDelete(mymsg.ChannelID, mymsg.ID)
		return
	},
}
