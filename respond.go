package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Respond responds (via DM if viaDM is set to true) with the text or embed provided. If both provided - only text is
// responded with.
func (req *Request) Respond(text string, embed *discordgo.MessageEmbed, viaDM bool) (m *discordgo.Message, err error) {
	// Determine which channel should we send the response to.
	var channelID string
	if viaDM {
		if channel, err := req.Sugo.Session.UserChannelCreate(req.Message.Author.ID); err != nil {
			return nil, err
		} else {
			channelID = channel.ID
		}
	} else {
		channelID = req.Message.ChannelID
	}

	// If text is provided - send text and nothing else.
	if text != "" {
		return req.Sugo.Session.ChannelMessageSend(channelID, text)
	}

	// If embed is provided - send embed.
	if embed != nil {
		return req.Sugo.Session.ChannelMessageSendEmbed(channelID, embed)
	}

	// Report error if no text and no embed are available for sending.
	return nil, errors.New("respond error: neither text nor embed were provided for sending")
}
