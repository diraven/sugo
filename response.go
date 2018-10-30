package sugo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type Response struct {
	Type    responseType
	Request *Request
	Text    string
	Embed   *discordgo.MessageEmbed
}

type responseType string

const (
	ResponsePlainText responseType = "plain_text"
	ResponseDefault   responseType = "default"
	ResponseInfo      responseType = "info"
	ResponseSuccess   responseType = "success"
	ResponseWarning   responseType = "warning"
	ResponseDanger    responseType = "danger"
)

func (req *Request) SimpleResponse(text string) (resp *Response) {
	return req.NewResponse(ResponseDefault, "", text)
}

// Respond responds (via DM if viaDM is set to true) with the Text or Embed provided. If both provided - only Text is
// responded with.
func (req *Request) NewResponse(respType responseType, title string, text string) (resp *Response) {
	resp = &Response{
		Type: respType,
	}

	resp.Request = req

	// If Response type is plain Text:
	if respType == ResponsePlainText {
		// Fill out the necessary fields and return the plain Text Response.
		resp.Text = text
		if title != "" {
			resp.Text = ">\n**" + title + "**\n" + resp.Text
		}
		return
	}

	// If title is not provided - generate it.
	if title == "" {
		title = "@" + req.Message.Author.Username
	}

	// If Response is not plain Text, create a basic Embed.
	resp.Embed = &discordgo.MessageEmbed{
		Title:       title,
		Description: text,
	}

	// Adjust Embed style depending on the message.
	switch respType {
	case ResponseDefault:
		resp.Embed.Color = ColorDefault
	case ResponseInfo:
		resp.Embed.Color = ColorInfo
		resp.Embed.Title = ":information_source:" + " " + resp.Embed.Title
	case ResponseSuccess:
		resp.Embed.Color = ColorSuccess
		resp.Embed.Title = ":white_check_mark:" + " " + resp.Embed.Title
	case ResponseWarning:
		resp.Embed.Color = ColorWarning
		resp.Embed.Title = ":warning:" + " " + resp.Embed.Title
	case ResponseDanger:
		resp.Embed.Color = ColorDanger
		resp.Embed.Title = ":exclamation:" + " " + resp.Embed.Title
	}

	// Return resulting Response.
	return
}

func (resp *Response) send(channelID string) (m *discordgo.Message, err error) {
	// Make sure valid Request is provided.
	if resp.Request == nil {
		return nil, errors.New("unable to send Response: empty Request provided")
	}

	switch resp.Type {
	case ResponsePlainText:
		// Response is a plain text response, send is a plain text.
		return resp.Request.Sugo.Session.ChannelMessageSend(channelID, resp.Text)

	case ResponseDefault, ResponseInfo, ResponseSuccess, ResponseWarning, ResponseDanger:
		// If response if one of the embed types - send response as an embed.
		return resp.Request.Sugo.Session.ChannelMessageSendEmbed(channelID, resp.Embed)
	}

	// Report error if it's some kind of weird unknown response type.
	return nil, errors.New("unknown response type")
}

// Send sends a Response into the channel Request was sent in.
func (resp *Response) Send() (m *discordgo.Message, err error) {
	return resp.send(resp.Request.Message.ChannelID)
}

// SendDM sends a Response to the user DirectMessages channel.
func (resp *Response) SendDM() (m *discordgo.Message, err error) {
	var channel *discordgo.Channel
	if channel, err = resp.Request.Sugo.Session.UserChannelCreate(resp.Request.Message.Author.ID); err != nil {
		return
	}
	return resp.send(channel.ID)
}
