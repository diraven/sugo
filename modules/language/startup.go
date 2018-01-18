package language

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Tables.
	if _, err := sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS ` + settingsTableName + ` (
			object_id TEXT,
			language TEXT
		);
	`); err != nil {
		return err
	}

	// Indexes.
	if _, err := sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS object_id ON ` + settingsTableName + ` (object_id);
	`); err != nil {
		return err
	}

	// Now load the data.
	return languages.reload(sg)
}
