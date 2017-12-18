package stats

import (
	"github.com/diraven/sugo"
	"time"
	"github.com/bwmarrin/discordgo"
)

type tStats struct{}

// logPlaying logs user status update and game he/she played if any.
func (s *tStats) logPlaying(sg *sugo.Instance, guildID string, userID string, theType int, game string) error {
	// Add new feed.
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO stats_playing (
			guild_id, user_id, the_type, game, created_at
		) VALUES (
			?, ?, ?, ?, ?
		);
	`, guildID, userID, theType, game, time.Now().Format("2006-01-02"))

	if err != nil {
		return err
	}

	// TODO Consider adding automatic clean up (such as remove all data records older then a month).

	return nil
}

func (s *tStats) getMostPlayedGames(sg *sugo.Instance, guildID string) ([]string, error) {
	// Variable to hold results.
	var gamesNames []string

	// Get rows from DB.
	rows, err := sg.DB.Query(`
		SELECT game, COUNT(*) as "count"
		FROM stats_playing
		WHERE DATE(created_at) > DATE('now', '-1 month')
		GROUP BY game
		ORDER BY count
		  DESC
		LIMIT 10;
	`)
	if err != nil {
		return nil, err
	}

	// Collect our games.
	for rows.Next() {
		var gameName string
		var count int
		if err := rows.Scan(&gameName, &count); err != nil {
			return nil, err
		}
		gamesNames = append(gamesNames, gameName)
	}

	return gamesNames, nil
}

// logPlaying logs user status update and game he/she played if any.
func (s *tStats) logMessage(sg *sugo.Instance, guildID string, userID string) error {
	// Add new feed.
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO stats_messaging (
			guild_id, user_id, created_at
		) VALUES (
			?, ?, ?
		);
	`, guildID, userID, time.Now().Format("2006-01-02 15:04:05"))

	if err != nil {
		return err
	}

	// TODO Consider adding automatic clean up (such as remove all data records older then a month).

	return nil
}

func (s *tStats) getMostMessagingUsers(sg *sugo.Instance, guildID string) ([]*discordgo.User, error) {
	// Variable to hold results.
	var users []*discordgo.User

	// Get rows from DB.
	rows, err := sg.DB.Query(`
		SELECT user_id, COUNT(*) as "count"
		FROM stats_messaging
		WHERE DATE(created_at) > DATE('now', '-1 month') AND guild_id=?
		GROUP BY user_id
		ORDER BY count
		  DESC
		LIMIT 10;
	`, guildID)
	if err != nil {
		return nil, err
	}

	// Collect the data into the usable format.
	for rows.Next() {
		var userID string
		var count int
		if err := rows.Scan(&userID, &count); err != nil {
			return nil, err
		}
		user, err := sg.User(userID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

//// add makes role public.
//func (s *tStats) updatePostedAt(sg *sugo.Instance, channelID string, url string, datetime *time.Time) error {
//	// Add new feed.
//	_, err := sg.DB.Exec(`
//		UPDATE feeds SET posted_at=?
//			WHERE channel_id=? AND url=?
//			;
//	`, datetime.Unix(), channelID, url)
//
//	if err != nil {
//		return err
//	}
//
//	err = f.reload(sg)
//	return err
//}
//
//// del makes rule not public.
//func (s *tStats) del(sg *sugo.Instance, channelID string, url string) error {
//	// Try to find a role in a list of public roles.
//	for _, item := range *f {
//		if item.ChannelID == channelID && item.Url == url {
//			// Perform deletion.
//			_, err := sg.DB.Exec(`
//				DELETE FROM feeds
//					WHERE channel_id=? AND url=?;
//			`, channelID, url)
//			if err != nil {
//				return err
//			}
//
//			// Reload data.
//			f.reload(sg)
//			return nil
//		}
//	}
//
//	// If feed not found, just notify user.
//	return errors.New("feed not found")
//}
//
//// reload reloads in-memory feeds list from the database.
//func (s *tStats) reload(sg *sugo.Instance) error {
//	// Variable to store errors if any.
//	var err error
//
//	// Initialize model storage.
//	*f = tStats{}
//
//	// Get rows from DB.
//	rows, err := sg.DB.Query("SELECT channel_id, url, posted_at FROM feeds")
//	if err != nil {
//		return err
//	}
//
//	// Put rows into the in-memory struct.
//	for rows.Next() {
//		var channelID string
//		var url string
//		var postedAt int64
//		if err := rows.Scan(&channelID, &url, &postedAt); err != nil {
//			return err
//		}
//		var postedAtTime = time.Unix(postedAt, 0)
//		*f = append(*f, tFeed{
//			channelID,
//			url,
//			&postedAtTime,
//		})
//	}
//
//	return nil
//}
