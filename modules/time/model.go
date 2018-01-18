package time

import (
	"github.com/diraven/sugo"
)

type tTimezones map[string]string

// set sets timezone setting for given object_id.
func (tz *tTimezones) set(sg *sugo.Instance, objectID string, timezone string) error {
	// Set timezone for object.
	if _, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO time_zones (
			object_id, timezone
		) VALUES (
			?, ?
		);
	`, objectID, timezone); err != nil {
		return err
	}

	// Reload data.
	if err := tz.reload(sg); err != nil {
		return err
	}

	return nil
}

// reset removes timezone setting.
func (tz *tTimezones) reset(sg *sugo.Instance, objectID string) error {
	// Set timezone for object.
	if _, err := sg.DB.Exec(`
		DELETE FROM time_zones
		WHERE object_id=?;
	`, objectID); err != nil {
		return err
	}

	// Reload data.
	if err := tz.reload(sg); err != nil {
		return err
	}

	return nil
}

// get returns timezone name for given object id.
func (tz *tTimezones) get(sg *sugo.Instance, req *sugo.Request) (string, error) {
	var objectID string

	// Get user ID.
	objectID = req.Message.Author.ID

	// Try to get user timezone.
	if timezone, ok := (*tz)[objectID]; ok {
		return timezone, nil
	}

	// Try to get guild timezone.
	if timezone, ok := (*tz)[req.Guild.ID]; ok {
		return timezone, nil
	}

	// No settings found, return default.
	return "UTC", nil
}

// reload reloads data from the database.
func (tz *tTimezones) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Initialize model storage.
	*tz = tTimezones{}

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT object_id, timezone FROM time_zones")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var objectID string
		var timezone string
		if err := rows.Scan(&objectID, &timezone); err != nil {
			return err
		}
		(*tz)[objectID] = timezone
	}

	return nil
}
