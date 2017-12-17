package sugo

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// RespondBadCommandUsage responds to the channel with "incorrect command usage" message mentioning person that invoked
// command.
func (sg *Instance) RespondBadCommandUsage(m *discordgo.Message, c *Command, title string, description string) (*discordgo.Message, error) {
	if title == "" {
		title = "bad command usage"
	}
	if description == "" {
		description = "see \"" + c.FullHelpPath(sg) + "\" for details"
	}
	msg, err := sg.RespondDanger(m, title, description)
	return msg, err
}

// RespondNotImplemented responds to the channel with "not implemented" message mentioning person that invoked
// command.
func (sg *Instance) RespondNotImplemented(m *discordgo.Message) (*discordgo.Message, error) {
	msg, err := sg.RespondWarning(m, "not implemented", "this functionality is not implemented yet")
	return msg, err
}

// RespondCommandNotFound responds to the channel with "command not found" message mentioning person that invoked
// command.
func (sg *Instance) RespondCommandNotFound(m *discordgo.Message) (*discordgo.Message, error) {
	return sg.RespondDanger(m, "command not found", "")
}

// Respond responds to the channel with an embed without any icons.
func (sg *Instance) Respond(m *discordgo.Message, title string, description string, color int, icon string) (*discordgo.Message, error) {
	if title == "" {
		title = "@" + m.Author.Username
	}
	if color == 0 {
		color = ColorDefault
	}
	msg, err := sg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       strings.Join([]string{icon, title}, " "),
		Description: description,
		Color:       color,
	})
	return msg, err
}

// RespondInfo responds to the channel with the "info" embed.
func (sg *Instance) RespondInfo(m *discordgo.Message, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(m, title, description, ColorInfo, ":information_source:")
}

// RespondInfo responds to the channel with the "success" embed.
func (sg *Instance) RespondSuccess(m *discordgo.Message, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(m, title, description, ColorSuccess, ":white_check_mark:")
}

// RespondInfo responds to the channel with the "warning" embed.
func (sg *Instance) RespondWarning(m *discordgo.Message, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(m, title, description, ColorWarning, ":warning:")
}

// RespondInfo responds to the channel with the "Danger" embed.
func (sg *Instance) RespondDanger(m *discordgo.Message, title string, description string) (*discordgo.Message, error) {
	return sg.Respond(m, title, description, ColorDanger, ":no_entry:")
}
