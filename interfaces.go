package sugo

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Match(sg *Instance, m *discordgo.Message) (matched bool, err error)
	CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error)
	Execute(sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed)
	Trigger() (trigger string)
	SetTrigger(trigger string)
}
