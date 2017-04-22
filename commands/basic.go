package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/errors"
)

type Basic struct {
	RootOnly            bool
	Triggers            []string
	PermissionsRequired []int
	Response            *string
}

func (c Basic) IsApplicable(sg *sugo.Instance, m *discordgo.Message) (is_applicable bool, err error) {
	is_applicable = false
	for _, trigger := range c.Triggers {
		if strings.Contains(m.Content, trigger) {
			is_applicable = true
			return
		}
	}
	return
}

func (c Basic) IsAllowed(sg *sugo.Instance, m *discordgo.Message) (result bool, err error) {
	// By default command is not allowed.
	result = false

	// For security reasons - every command should have at least one permission set explicitly.
	if len(c.PermissionsRequired) == 0 {
		err = errors.SugoError{Text: "Command has no PermissionsRequired[]!"}
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

	// If root user issued a command - it is always allowed.
	if sg.Root != nil {
		if m.Author.ID == sg.Root.ID {
			return true, nil
		}
	}
	if c.RootOnly {
		return // If command is for root only - we do not check anything else and just deny using it.
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
	result = true
	return
}

func (c Basic) Execute(sg *sugo.Instance, m *discordgo.Message) (err error) {
	_, err = sg.ChannelMessageSend(m.ChannelID, *c.Response)
	return
}
