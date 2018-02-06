package sugo

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// RespondBadCommandUsage responds to the channel with "incorrect command usage" message mentioning person that invoked
// command.
func (req *Request) RespondBadCommandUsage(title string, description string) (*discordgo.Message, error) {
	if title == "" {
		title = "bad command usage"
	}
	if description == "" {
		description = "command used incorrectly"
	}
	msg, err := req.RespondDanger(title, description)
	return msg, err
}

// RespondGuildOnly responds to the channel with "guild only" message explaining that command is only available in guild
// channel.
func (req *Request) RespondGuildOnly(title string, description string) (*discordgo.Message, error) {
	if title == "" {
		title = "guild only"
	}
	if description == "" {
		description = "this command can only be executed in guild channel"
	}
	msg, err := req.RespondDanger(title, description)
	return msg, err
}

// RespondNotImplemented responds to the channel with "not implemented" message mentioning person that invoked
// command.
func (req *Request) RespondNotImplemented() (*discordgo.Message, error) {
	msg, err := req.RespondWarning("not implemented", "this functionality is not implemented yet")
	return msg, err
}

// RespondCommandNotFound responds to the channel with "command not found" message mentioning person that invoked
// command.
func (req *Request) RespondCommandNotFound() (*discordgo.Message, error) {
	return req.RespondWarning("command not found", "")
}

// Respond responds to the channel with an embed without any icons.
func (req *Request) Respond(title string, description string, color int, icon string) (*discordgo.Message, error) {
	if title == "" {
		title = "@" + req.Message.Author.Username
	}
	if color == 0 {
		color = ColorDefault
	}
	msg, err := req.Sugo.Session.ChannelMessageSendEmbed(req.Message.ChannelID, &discordgo.MessageEmbed{
		Title:       strings.Join([]string{icon, title}, " "),
		Description: description,
		Color:       color,
	})
	return msg, err
}

// RespondInfo responds to the channel with the "info" embed.
func (req *Request) RespondInfo(title string, description string) (*discordgo.Message, error) {
	return req.Respond(title, description, ColorInfo, ":information_source:")
}

// RespondSuccess responds to the channel with the "success" embed.
func (req *Request) RespondSuccess(title string, description string) (*discordgo.Message, error) {
	return req.Respond(title, description, ColorSuccess, ":white_check_mark:")
}

// RespondWarning responds to the channel with the "warning" embed.
func (req *Request) RespondWarning(title string, description string) (*discordgo.Message, error) {
	return req.Respond(title, description, ColorWarning, ":warning:")
}

// RespondDanger responds to the channel with the "Danger" embed.
func (req *Request) RespondDanger(title string, description string) (*discordgo.Message, error) {
	return req.Respond(title, description, ColorDanger, ":no_entry:")
}
