package sugo

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Validate(sg *Instance, m *discordgo.Message) (passed bool, err error)
	CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error)
	Execute(sg *Instance, m *discordgo.Message) (err error)
}
