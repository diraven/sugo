package sugo

import (
	"github.com/bwmarrin/discordgo"
	"errors"
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
func (r *Request) GetGuild() (*discordgo.Guild, error) {
	if r.Channel.GuildID != "" {
		guild, err := r.Sugo.Session.State.Guild(r.Channel.GuildID)
		if err != nil {
			return nil, errors.New("unable to get guild for request: " + r.Query)
		}

		return guild, nil
	}

	return nil, errors.New("request has no guild")
}

// IsChannelDefault returns true if channel is Guild's default channel and false otherwise.
func (r *Request) IsChannelDefault() bool {
	if r.Channel.ID == r.Channel.GuildID {
		return true
	}

	return false
}

// IsChannelDM returns true if channel is DirectMessages (or GroupDirectMessages) channel and false otherwise.
func (r *Request) IsChannelDM() bool {
	if r.Channel.Type == discordgo.ChannelTypeDM || r.Channel.Type == discordgo.ChannelTypeGroupDM {
		return true
	}

	return false
}

// WrapError error wraps error with additional request info.
func (r *Request) WrapError(e error, text string) error {
	return errors.New("request error: " + r.Command.GetPath() + ": " + r.Query)
}
