package time

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Tables.
	if _, err := sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS time_zones (
			object_id TEXT,
			timezone TEXT
		);
	`); err != nil {
		return err
	}

	// Indexes.
	if _, err := sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS object ON time_zones (object_id);
	`); err != nil {
		return err
	}

	// Now load the data.
	return timezones.reload(sg)
}
