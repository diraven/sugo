package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"fmt"
	"github.com/diraven/sugo/helpers"
)

type Basic struct {
	RootOnly            bool
	PermissionsRequired []int
	Response            string
}

func (c Basic) Validate(sg *sugo.Instance, m *discordgo.Message) (passed bool, err error) {
	passed = true
	return
}

func (c Basic) CheckPermissions(sg *sugo.Instance, m *discordgo.Message) (passed bool, err error) {
	// By default command is not allowed.
	passed = false

	// For security reasons - every command should have at least one permission set explicitly.
	if len(c.PermissionsRequired) == 0 {
		err = sugo.Error{Text: "Command has no PermissionsRequired[]!"}
		return
	}

	// Get channel to check for permissions.
	channel, err := sg.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Calculate compound permission.
	var compound_perm int = 0
	for _, perm := range c.PermissionsRequired {
		compound_perm |= perm
	}

	// Make sure bot has the permission required.
	bot_has_perm, err := sg.BotHasPermission(compound_perm, channel)
	if err != nil {
		return
	}
	if !(bot_has_perm) {
		return
	}

	// If user is a root - command is always allowed.
	if sg.IsRoot(m.Author) {
		return true, nil
	}
	// Otherwise if user is not a root a command is root-only - command is not allowed.
	if c.RootOnly {
		return
	}

	// Make sure user has the permission required.
	user_has_perm, err := sg.UserHasPermission(compound_perm, m.Author, channel)
	if err != nil {
		return
	}
	if !(user_has_perm) {
		return
	}

	// At this time we have checked that:
	// - Command has at least one permission requirement.
	// - Channel we check permissions against exists.
	// - User has all the permissions required.
	// - Bot has all the permissions required.
	// So we can safely say the command is allowed to be executed.
	passed = true
	return
}

func (c Basic) Execute(sg *sugo.Instance, m *discordgo.Message) (err error) {
	_, err = c.RespondMention(sg, m, c.Response)
	if err != nil {
		return
	}
	return
}

// Responds to the channel.
func (c Basic) Respond(sg *sugo.Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

// Responds to the channel with the original message author mention.
func (c Basic) RespondMention(sg *sugo.Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	response_text := fmt.Sprintf("%s %s", helpers.UserAsMention(m.Author), text)
	message, err = sg.ChannelMessageSend(m.ChannelID, response_text)
	if err != nil {
		return
	}
	return
}
