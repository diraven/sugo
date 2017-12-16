package public_roles

import (
	"github.com/diraven/sugo"
)

func startup(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Prepare database structure to store the data.
	// Public Roles table.
	_, err = sg.DB.Exec(`
		CREATE TABLE IF NOT EXISTS public_roles (
			guild_id TEXT,
			role_id TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Index to enforce guild_id and role_id uniqueness.
	_, err = sg.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS guild_role ON public_roles (guild_id, role_id);
	`)
	if err != nil {
		return err
	}

	// Now load the data.
	return publicRoles.reload(sg)
}
