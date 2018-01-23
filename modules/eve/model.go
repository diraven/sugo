package eve

import (
	"github.com/diraven/sugo"
)

type tKillmails struct {
	allianceIDs    map[string]string
	corporationIDs map[string]string
}

func (s *tKillmails) addAlliance(sg *sugo.Instance, channelID string, eveID string) error {
	// Add new feed.
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO eve_killmail_alliances (
			channel_id, eve_id
		) VALUES (
			?, ?
		);
	`, channelID, eveID)

	if err != nil {
		return err
	}

	err = s.reload(sg)
	return err
}

func (s *tKillmails) addCorporation(sg *sugo.Instance, channelID string, eveID string) error {
	// Add new feed.
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO eve_killmail_corporations (
			channel_id, eve_id
		) VALUES (
			?, ?
		);
	`, channelID, eveID)

	if err != nil {
		return err
	}

	err = s.reload(sg)
	return err
}

func (s *tKillmails) delAlliance(sg *sugo.Instance, channelID string, eveID string) error {
	// Perform deletion.
	_, err := sg.DB.Exec(`
		DELETE FROM eve_killmail_alliances
			WHERE channel_id=? AND eve_id=?;
	`, channelID, eveID)
	if err != nil {
		return err
	}

	// Reload data.
	s.reload(sg)
	return nil
}

func (s *tKillmails) delCorporation(sg *sugo.Instance, channelID string, eveID string) error {
	// Perform deletion.
	_, err := sg.DB.Exec(`
		DELETE FROM eve_killmail_corporations
			WHERE channel_id=? AND eve_id=?;
	`, channelID, eveID)
	if err != nil {
		return err
	}

	// Reload data.
	s.reload(sg)
	return nil
}

// reload reloads in-memory feeds list from the database.
func (s *tKillmails) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Initialize model storage.
	*s = tKillmails{}

	// Initalize data storages.
	s.corporationIDs = map[string]string{}
	s.allianceIDs = map[string]string{}

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT channel_id, eve_id FROM eve_killmail_alliances")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var channelID string
		var eveID string
		if err := rows.Scan(&channelID, &eveID); err != nil {
			return err
		}
		s.allianceIDs[eveID] = channelID
	}

	// Get rows from DB.
	rows, err = sg.DB.Query("SELECT channel_id, eve_id FROM eve_killmail_corporations")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var channelID string
		var eveID string
		if err := rows.Scan(&channelID, &eveID); err != nil {
			return err
		}
		s.corporationIDs[eveID] = channelID
	}

	return nil
}
