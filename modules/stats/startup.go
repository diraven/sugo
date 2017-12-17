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
			created_at TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce guild, user, game and created_at uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS guild_user_game_created_at ON stats_playing (guild_id, user_id, game, created_at);
	`)
	if err != nil {
		return err
	}

	return nil
}
