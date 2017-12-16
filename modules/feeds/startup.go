package feeds

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Public Roles table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS feeds (
			channel_id TEXT,
			url TEXT,
			posted_at INTEGER
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce channel_id and url uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS channel_url ON feeds (channel_id, url);
	`)
	if err != nil {
		return err
	}

	// Start posting items.
	go postNewItems(sg)

	// Now load the data.
	return feeds.reload(sg)
}
