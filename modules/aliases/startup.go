package aliases

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Aliases table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS aliases (
			guild_id TEXT,
			alias TEXT,
			command_path TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce guild_id and alias uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS guild_alias ON aliases (guild_id, alias);
	`)
	if err != nil {
		return err
	}

	// Now load the data.
	return aliases.reload(sg)
}
