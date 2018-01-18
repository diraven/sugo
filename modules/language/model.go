package language

import (
	"github.com/diraven/sugo"
)

type tLanguages map[string]string

const (
	settingsTableName = "language_settings"
)

// set sets language setting for given object_id.
func (tz *tLanguages) set(sg *sugo.Instance, objectID string, language string) error {
	// Set language for object.
	if _, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO `+settingsTableName+` (
			object_id, language
		) VALUES (
			?, ?
		);
	`, objectID, language); err != nil {
		return err
	}

	// Reload data.
	if err := tz.reload(sg); err != nil {
		return err
	}

	return nil
}

// reset removes language setting.
func (tz *tLanguages) reset(sg *sugo.Instance, objectID string) error {
	// Set language for object.
	if _, err := sg.DB.Exec(`
		DELETE FROM `+settingsTableName+`
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

// get returns language name for given object id.
func (tz *tLanguages) get(sg *sugo.Instance, req *sugo.Request) (string, error) {
	var objectID string

	// Get user ID.
	objectID = req.Message.Author.ID

	// Try to get user language.
	if language, ok := (*tz)[objectID]; ok {
		return language, nil
	}

	// Try to get guild language.
	if language, ok := (*tz)[req.Guild.ID]; ok {
		return language, nil
	}

	// No settings found, return default.
	return defaultLanguage, nil
}

// reload reloads data from the database.
func (tz *tLanguages) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Initialize model storage.
	*tz = tLanguages{}

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT object_id, language FROM " + settingsTableName)
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var objectID string
		var language string
		if err := rows.Scan(&objectID, &language); err != nil {
			return err
		}
		(*tz)[objectID] = language
	}

	return nil
}
