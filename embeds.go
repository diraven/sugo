package sugo

import (
	"github.com/bwmarrin/discordgo"
)

func NewDefaultEmbed(req *Request, desc string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "@" + req.Message.Author.Username,
		Description: desc,
		Color:       ColorDefault,
	}
}

// NewInfoEmbed creates info-styled embed.
func NewInfoEmbed(req *Request, desc string) (embed *discordgo.MessageEmbed) {
	embed = NewDefaultEmbed(req, desc)
	embed.Color = ColorInfo
	embed.Title = ":information_source:" + " " + embed.Title
	return embed
}

// NewSuccessEmbed creates success-styled embed.
func NewSuccessEmbed(req *Request, desc string) (embed *discordgo.MessageEmbed) {
	embed = NewDefaultEmbed(req, desc)
	embed.Color = ColorSuccess
	embed.Title = ":white_check_mark:" + " " + embed.Title
	return embed
}

// NewWarningEmbed creates warning-styled embed.
func NewWarningEmbed(req *Request, desc string) (embed *discordgo.MessageEmbed) {
	embed = NewDefaultEmbed(req, desc)
	embed.Color = ColorWarning
	embed.Title = ":warning:" + " " + embed.Title
	return embed
}

// NewDangerEmbed creates danger-styled embed.
func NewDangerEmbed(req *Request, desc string) (embed *discordgo.MessageEmbed) {
	embed = NewDefaultEmbed(req, desc)
	embed.Color = ColorDanger
	embed.Title = ":exclamation:" + " " + embed.Title
	return embed
}
