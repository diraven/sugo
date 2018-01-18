package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/i18n"
)

// Request contains all given request related info.
type Request struct {
	Message       *discordgo.Message
	Command       *Command
	Query         string
	Channel       *discordgo.Channel
	Guild         *discordgo.Guild
	TranslateFunc *i18n.TranslateFunc
	IsCanceled    bool
}
