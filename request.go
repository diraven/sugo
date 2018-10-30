package sugo

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Request contains message context data along with some helpers to retrieve more information.
type Request struct {
	Ctx     context.Context
	Sugo    *Instance
	Message *discordgo.Message
	Channel *discordgo.Channel
	Command *Command
	Query   string
}

// GetGuild allows to retrieve *discordgo.Guild from Request. Will not work and will throw error for channels
// that have no guild such as DirectMessages or GroupDirectMessages channels, so you probably want to check
// those beforehand.
func (req *Request) GetGuild() (*discordgo.Guild, error) {
	if req.Channel.GuildID != "" {
		guild, err := req.Sugo.Session.State.Guild(req.Channel.GuildID)
		if err != nil {
			return nil, errors.New("unable to get guild for Request")
		}

		return guild, nil
	}

	return nil, errors.New("Request has no guild")
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
