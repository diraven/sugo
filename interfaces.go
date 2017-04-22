package sugo

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	IsApplicable(sg *Instance, m *discordgo.Message) (is_applicable bool, err error)
	IsAllowed(sg *Instance, m *discordgo.Message) (passed bool, err error)
	Execute(sg *Instance, m *discordgo.Message) (err error)
}
