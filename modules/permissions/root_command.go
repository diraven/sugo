package permissions

import (
	"github.com/diraven/sugo"
)

var rootCommand = &sugo.Command{
	Trigger:     "permissions",
	RootOnly:    true,
	Description: "Allows to manipulate custom command permissions.",
	SubCommands: []*sugo.Command{
		{
			Trigger:     "allow",
			Description: "Allows specific command usage for specific role.",
			Usage:       "@role command [subcommand ...]",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var err error

				// Try to find role.
				role, q := findRole(sg, req, req.Query)
				if role == nil {
					// Role not found.
					sg.RespondDanger(req, "", "role not found")
					return nil
				}

				// Try to find command.
				command, err := sg.FindCommand(req, q)
				if command == nil {
					// Command not found.
					sg.RespondCommandNotFound(req)
					return nil
				}

				// Both command and role exist.
				err = permissions.set(sg, role.ID, command.Path(), true)
				if err != nil {
					return err
				}

				// Notify user of success.
				sg.RespondSuccess(req, "", "")
				return nil
			},
		},
		{
			Trigger:     "deny",
			Description: "Denies specific command usage for specific role.",
			Usage:       "@role command [subcommand ...]",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var err error

				// Try to find role.
				role, q := findRole(sg, req, req.Query)
				if role == nil {
					sg.RespondDanger(req, "", "role not found")
					return err
				}

				// Try to find command.
				command, err := sg.FindCommand(req, q)
				if command == nil {
					sg.RespondCommandNotFound(req)
					return err
				}

				permissions.set(sg, role.ID, command.Path(), false)
				sg.RespondSuccess(req, "", "")
				return err
			},
		},
		{
			Trigger:     "default",
			Description: "Sets permissions for specific command usage by role to default.",
			Usage:       "@role command [subcommand ...]",
			AllowParams: true,
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var err error

				// Try to find role.
				role, q := findRole(sg, req, req.Query)
				if role == nil {
					sg.RespondDanger(req, "", "role not found")
					return err
				}

				// Try to find command.
				command, err := sg.FindCommand(req, q)
				if command == nil {
					sg.RespondCommandNotFound(req)
					return err
				}

				permissions.setDefault(sg, role.ID, command.Path())
				sg.RespondSuccess(req, "", "")
				return err
			},
		},
		{
			Trigger:     "roles",
			Description: "Shows all server roles with their IDs.",
			Execute: func(sg *sugo.Instance, req *sugo.Request) error {
				var err error

				// Response text.
				response := "```\n"

				// For each guild role.
				for _, role := range req.Guild.Roles {
					// If response is too long already - make a new one.
					if len(response) > 1500 {
						response = response + "```"
						_, err = sg.RespondInfo(req, "", response)
						response = "```\n"
					}
					response = response + role.ID + ": " + role.Name + "\n"
				}

				// End response text.
				response = response + "```"

				_, err = sg.RespondInfo(req, "", response)
				return err
			},
		},
	},
}
