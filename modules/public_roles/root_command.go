package public_roles

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"context"
)

var rootCommand = &sugo.Command{
	Trigger:            "public_roles",
	Description:        "Allows to manipulate public roles.",
	PermittedByDefault: true,
	ParamsAllowed:      true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, m, q)

		// Start building response.
		var response string

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response = response + "```\n"
			response = response + sugo.FmtStringsSlice(rolesToRoleNames(roles), "\n", 1500, "\n...", "")
			response = response + "```"
			_, err = sg.RespondInfo(m, "public roles", response)
		} else {
			_, err = sg.RespondDanger(m, "", "nothing found")
		}

		return err
	},
	SubCommands: []*sugo.Command{
		myCmd,
		whoCmd,
		addCmd,
		delCmd,
		joinCmd,
		leaveCmd,
		createCmd,
		statsCmd,
	},
}
