package sugo

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

// Request contains message context data along with some helpers to retrieve more information.
type Request struct {
	Sugo    *Instance
	Message *discordgo.Message
	Channel *discordgo.Channel
	Command *Command
	Query   string
}

// GetGuild allows to retrieve *discordgo.Guild from request. Will not work and will throw error for channels
// that have no guild such as DirectMessages or GroupDirectMessages channels, so you probably want to check
// those beforehand.
func (req *Request) GetGuild() (*discordgo.Guild, error) {
	if req.Channel.GuildID != "" {
		guild, err := req.Sugo.Session.State.Guild(req.Channel.GuildID)
		if err != nil {
			return nil, errors.New("unable to get guild for request: " + req.Query)
		}

		return guild, nil
	}

	return nil, errors.New("request has no guild")
}

// IsChannelDefault returns true if channel is Guild's default channel and false otherwise.
func (req *Request) IsChannelDefault() bool {
	if req.Channel.ID == req.Channel.GuildID {
		return true
	}

	return false
}

// IsChannelDM returns true if channel is DirectMessages (or GroupDirectMessages) channel and false otherwise.
func (req *Request) IsChannelDM() bool {
	if req.Channel.Type == discordgo.ChannelTypeDM || req.Channel.Type == discordgo.ChannelTypeGroupDM {
		return true
	}

	return false
}

// WrapError error wraps error with additional request info.
func (req *Request) WrapError(e error, text string) error {
	return errors.New("request error: " + req.Command.GetPath() + ": " + req.Query)
}
