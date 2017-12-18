package stats

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Table.
	if _, err := sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS stats_playing (
			guild_id TEXT,
			user_id TEXT,
			the_type TINYINT,
			game TEXT,
			created_at TEXT
		);
	`); err != nil {
		return err
	}

	// Indexes.
	if _, err := sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS guild_user_game_created_at ON stats_playing (guild_id, user_id, game, created_at);
	`); err != nil {
		return err
	}

	// Table.
	if _, err := sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS stats_messaging (
			guild_id TEXT,
			user_id TEXT,
			created_at TEXT
		);
	`); err != nil {
		return err
	}

	//// Indexes.
	//if _, err := sg.DB.Exec(`
	//	CREATE UNIQUE INDEX IF NOT EXISTS guild_user_game_created_at ON stats_playing (guild_id, user_id, game, created_at);
	//`); err != nil {
	//	return err
	//}

	return nil
}
