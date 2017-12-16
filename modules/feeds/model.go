package feeds

import (
	"github.com/diraven/sugo"
	"errors"
	"time"
)

type tFeed struct {
	ChannelID string
	Url       string
	PostedAt  *time.Time
}

type tFeeds []tFeed

// add makes role public.
func (f *tFeeds) add(sg *sugo.Instance, channelID string, url string) error {
	// Add new feed.
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO feeds (
			channel_id, url, posted_at
		) VALUES (
			?, ?, ?
		);
	`, channelID, url, time.Now().Unix())

	if err != nil {
		return err
	}

	err = f.reload(sg)
	return err
}

// add makes role public.
func (f *tFeeds) updatePostedAt(sg *sugo.Instance, channelID string, url string, datetime *time.Time) error {
	// Add new feed.
	_, err := sg.DB.Exec(`
		UPDATE feeds SET posted_at=?
			WHERE channel_id=? AND url=?
			;
	`, datetime.Unix(), channelID, url)

	if err != nil {
		return err
	}

	err = f.reload(sg)
	return err
}

// del makes rule not public.
func (f *tFeeds) del(sg *sugo.Instance, channelID string, url string) error {
	// Try to find a role in a list of public roles.
	for _, item := range *f {
		if item.ChannelID == channelID && item.Url == url {
			// Perform deletion.
			_, err := sg.DB.Exec(`
				DELETE FROM feeds
					WHERE channel_id=? AND url=?;
			`, channelID, url)
			if err != nil {
				return err
			}

			// Reload data.
			f.reload(sg)
			return nil
		}
	}

	// If feed not found, just notify user.
	return errors.New("feed not found")
}

// reload reloads in-memory feeds list from the database.
func (f *tFeeds) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Initialize model storage.
	*f = tFeeds{}

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT channel_id, url, posted_at FROM feeds")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var channelID string
		var url string
		var postedAt int64
		if err := rows.Scan(&channelID, &url, &postedAt); err != nil {
			return err
		}
		var postedAtTime = time.Unix(postedAt, 0)
		*f = append(*f, tFeed{
			channelID,
			url,
			&postedAtTime,
		})
	}

	return nil
}
