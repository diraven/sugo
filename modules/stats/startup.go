package stats

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Public Roles table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS stats_playing (
			guild_id TEXT,
			user_id TEXT,
			the_type TINYINT,
			game TEXT,
			created_at INTEGER
		);
	`)
	if err != nil {
		return err
	}

	//// Index to enforce channel_id and url uniqueness.
	//_, err = sg.DB.Exec(`
	//	CREATE UNIQUE INDEX IF NOT EXISTS channel_url ON feeds (channel_id, url);
	//`)
	//if err != nil {
	//	return err
	//}

	return nil
}
