package permissions

import (
	"github.com/diraven/sugo"
)

type tPermissionsStorage struct {
	Permissions map[string]bool
}

// get returns given permission status.
func (p *tPermissionsStorage) get(sg *sugo.Instance, roleID string, commandPath string) (isAllowed bool, exists bool) {
	var key string
	key = roleID + ":" + commandPath
	isAllowed, ok := p.Permissions[key]
	if ok {
		return isAllowed, true
	} else {
		return isAllowed, false
	}
}

// set sets given permission status.
func (p *tPermissionsStorage) set(sg *sugo.Instance, roleID string, commandPath string, isAllowed bool) error {
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO permissions (
			role_id, command_path, is_allowed
		) VALUES (
			?, ?, ?
		);
	`, roleID, commandPath, isAllowed)

	if err != nil {
		return err
	}

	err = p.reload(sg)
	return err
}

// setDefault resets given permission status to default state.
func (p *tPermissionsStorage) setDefault(sg *sugo.Instance, roleID string, commandPath string) error {
	_, err := sg.DB.Exec(`
		DELETE FROM permissions
			WHERE role_id=? AND command_path=?;
	`, roleID, commandPath)

	if err != nil {
		return err
	}

	err = p.reload(sg)
	return err
}

// reload reloads in-memory permissions cache from the database.
func (p *tPermissionsStorage) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Now load everything from the database into memory.
	p.Permissions = make(map[string]bool)

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT role_id, command_path, is_allowed FROM permissions")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var roleID string
		var commandPath string
		var isAllowed bool
		if err := rows.Scan(&roleID, &commandPath, &isAllowed); err != nil {
			return err
		}
		p.Permissions[roleID+":"+commandPath] = isAllowed
	}

	return err
}
