package permissions

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
)

var rootCommand = &sugo.Command{
	Trigger:     "permissions",
	RootOnly:    true,
	Description: "Allows to manipulate custom command permissions.",
	SubCommands: []*sugo.Command{
		{
			Trigger:       "allow",
			Description:   "Allows specific command usage for specific role.",
			Usage:         "@role command [subcommand ...]",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				// Try to find role.
				role, q := findRole(sg, m, q)
				if role == nil {
					// Role not found.
					sg.RespondDanger(m, "", "role not found")
					return nil
				}

				// Try to find command.
				command, err := sg.FindCommand(m, q)
				if command == nil {
					// Command not found.
					sg.RespondCommandNotFound(m)
					return nil
				}

				// Both command and role exist.
				err = permissions.set(sg, role.ID, command.Path(), true)
				if err != nil {
					return err
				}

				// Notify user of success.
				sg.RespondSuccess(m, "", "")
				return nil
			},
		},
		{
			Trigger:       "deny",
			Description:   "Denies specific command usage for specific role.",
			Usage:         "@role command [subcommand ...]",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				// Try to find role.
				role, q := findRole(sg, m, q)
				if role == nil {
					sg.RespondDanger(m, "", "role not found")
					return err
				}

				// Try to find command.
				command, err := sg.FindCommand(m, q)
				if command == nil {
					sg.RespondCommandNotFound(m)
					return err
				}

				permissions.set(sg, role.ID, command.Path(), false)
				sg.RespondSuccess(m, "", "")
				return err
			},
		},
		{
			Trigger:       "default",
			Description:   "Sets permissions for specific command usage by role to default.",
			Usage:         "@role command [subcommand ...]",
			ParamsAllowed: true,
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				// Try to find role.
				role, q := findRole(sg, m, q)
				if role == nil {
					sg.RespondDanger(m, "", "role not found")
					return err
				}

				// Try to find command.
				command, err := sg.FindCommand(m, q)
				if command == nil {
					sg.RespondCommandNotFound(m)
					return err
				}

				permissions.setDefault(sg, role.ID, command.Path())
				sg.RespondSuccess(m, "", "")
				return err
			},
		},
		{
			Trigger:     "roles",
			Description: "Shows all server roles with their IDs.",
			Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
				var err error

				// Get guild.
				guild, err := sg.GuildFromMessage(m)
				if err != nil {
					return err
				}

				// Response text.
				response := "```\n"

				// For each guild role.
				for _, role := range guild.Roles {
					// If response is too long already - make a new one.
					if len(response) > 1500 {
						response = response + "```"
						_, err = sg.RespondInfo(m, "", response)
						response = "```\n"
					}
					response = response + role.ID + ": " + role.Name + "\n"
				}

				// End response text.
				response = response + "```"

				_, err = sg.RespondInfo(m, "", response)
				return err
			},
		},
	},
}
