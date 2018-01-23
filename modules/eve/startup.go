package eve

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Public Roles table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS eve_killmail_alliances (
			channel_id TEXT,
			eve_id TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce channel_id and url uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS channel_alliance ON eve_killmail_alliances (channel_id, eve_id);
	`)
	if err != nil {
		return err
	}

	// Prepare database structure to store the data.
	// Public Roles table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS eve_killmail_corporations (
			channel_id TEXT,
			eve_id TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce channel_id and url uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS channel_corporation ON eve_killmail_corporations (channel_id, eve_id);
	`)
	if err != nil {
		return err
	}

	// Start posting items.
	go postKillmails(sg)

	// Now load the data.
	return killmails.reload(sg)
}
