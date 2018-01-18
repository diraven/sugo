package public_roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

var rootCommand = &sugo.Command{
	Trigger:            "public_roles",
	Description:        "Allows to manipulate public roles.",
	PermittedByDefault: true,
	AllowParams:        true,
	Execute: func(sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		// Try to find role based on query.
		roles, err := publicRoles.findGuildPublicRole(sg, m, q)

		// If we have got at least one suggested role.
		if len(roles) > 0 {
			// Make an array of suggested role names.
			response := sugo.FmtStringsSlice(rolesToRoleNames(roles), ", ", "`", 1500, "...", "")
			_, err = sg.RespondInfo(m, "available public roles", response)
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
