package commands

import (
	"github.com/diraven/sugo"
)

type Composite struct {
	Basic
	Subcommands []*sugo.Command
}

//CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error)
//Execute(sg *Instance, m *discordgo.Message) (err error)
//HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed)
