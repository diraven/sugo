package triggers

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Triggers table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS triggers (
			guild_id TEXT,
			trigger TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce guild uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS guild ON triggers (guild_id);
	`)
	if err != nil {
		return err
	}

	// Now load the data.
	return triggers.reload(sg)
}
