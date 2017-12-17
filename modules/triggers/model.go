package triggers

import (
	"github.com/diraven/sugo"
)

type tTriggers map[string]string

// get returns trigger for given guild
func (ts *tTriggers) get(sg *sugo.Instance, guildID string) string {
	trigger, ok := (*ts)[guildID]
	if ok {
		return trigger
	}
	return ""
}

// sets trigger for given guild
func (ts *tTriggers) set(sg *sugo.Instance, guildID string, trigger string) error {
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO triggers (
			guild_id, trigger
		) VALUES (
			?, ?
		);
	`, guildID, trigger)

	if err != nil {
		return err
	}

	err = ts.reload(sg)
	return err
}

// setDefault resets given guild trigger to default.
func (ts *tTriggers) setDefault(sg *sugo.Instance, guildID string) error {
	_, err := sg.DB.Exec(`
		DELETE FROM triggers
			WHERE guild_id=?;
	`, guildID)

	if err != nil {
		return err
	}

	err = ts.reload(sg)
	return err
}

// reload reloads in-memory triggers cache from the database.
func (ts *tTriggers) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Now load everything from the database into memory.
	*ts = tTriggers{}

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT guild_id, trigger FROM triggers")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var guildID string
		var trigger string
		if err := rows.Scan(&guildID, &trigger); err != nil {
			return err
		}
		(*ts)[guildID] = trigger
	}

	return err
}
