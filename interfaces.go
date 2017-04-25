package sugo

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Match(sg *Instance, m *discordgo.Message) (matched bool, err error)
	CheckPermissions(sg *Instance, m *discordgo.Message) (passed bool, err error)
	Execute(sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed(sg *Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed)

	Trigger() (value string)
	SetTrigger(value string)

	RootOnly() (value bool)
	SetRootOnly(value bool)

	PermissionsRequired() (value []int)
	AddRequiredPermission(value int)

	Response() (value string)
	SetResponse(value string)

	EmbedResponse() (value *discordgo.MessageEmbed)
	SetEmbedResponse(value *discordgo.MessageEmbed)

	Description() (value string)
	SetDescription(value string)

	Usage() (value string)
	SetUsage(value string)
}
