package sugo

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// RespondBadCommandUsage responds to the channel with "incorrect command usage" message mentioning person that invoked
// command.
func (sg *Instance) RespondBadCommandUsage(req *Request, title string, description string) (*discordgo.Message, error) {
	if title == "" {
		title = "bad command usage"
	}
	if description == "" {
		description = "command used incorrectly"
	}
	msg, err := sg.RespondDanger(req, title, description)
	return msg, err
}

// RespondGuildOnly responds to the channel with "guild only" message explaining that command is only available in guild
// channel.
func (sg *Instance) RespondGuildOnly(req *Request, title string, description string) (*discordgo.Message, error) {
	if title == "" {
		title = "guild only"
	}
	if description == "" {
		description = "this command can only be executed in guild channel"
	}
	msg, err := sg.RespondDanger(req, title, description)
	return msg, err
}

// RespondNotImplemented responds to the channel with "not implemented" message mentioning person that invoked
// command.
func (sg *Instance) RespondNotImplemented(req *Request) (*discordgo.Message, error) {
	msg, err := sg.RespondWarning(req, "not implemented", "this functionality is not implemented yet")
	return msg, err
}

// RespondCommandNotFound responds to the channel with "command not found" message mentioning person that invoked
// command.
func (sg *Instance) RespondCommandNotFound(req *Request) (*discordgo.Message, error) {
	return sg.RespondWarning(req, "command not found", "")
}

// Respond responds to the channel with an embed without any icons.
func (sg *Instance) Respond(req *Request, title string, description string, color int, icon string) (*discordgo.Message, error) {
	if title == "" {
		title = "@" + req.Message.Author.Username
	}
	if color == 0 {
		color = ColorDefault
	}
	msg, err := sg.Session.ChannelMessageSendEmbed(req.Message.ChannelID, &discordgo.MessageEmbed{
		Title:       strings.Join([]string{icon, title}, " "),
		Description: description,
		Color:       color,
	})
	return msg, err
}

// RespondInfo responds to the channel with the "info" embed.
func (sg *Instance) RespondInfo(req *Request, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(req, title, description, ColorInfo, ":information_source:")
}

// RespondSuccess responds to the channel with the "success" embed.
func (sg *Instance) RespondSuccess(req *Request, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(req, title, description, ColorSuccess, ":white_check_mark:")
}

// RespondWarning responds to the channel with the "warning" embed.
func (sg *Instance) RespondWarning(req *Request, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(req, title, description, ColorWarning, ":warning:")
}

// RespondDanger responds to the channel with the "Danger" embed.
func (sg *Instance) RespondDanger(req *Request, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(req, title, description, ColorDanger, ":no_entry:")
}
