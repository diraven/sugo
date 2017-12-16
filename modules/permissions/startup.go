package permissions

import "github.com/diraven/sugo"

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Permissions table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS permissions (
			role_id TEXT,
			command_path TEXT,
			is_allowed BOOLEAN
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce role_id and command uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS role_command ON permissions (role_id, command_path);
	`)
	if err != nil {
		return err
	}

	// Now load the data.
	return permissions.reload(sg)
}
