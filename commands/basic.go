package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"fmt"
	"github.com/diraven/sugo/helpers"
	"strings"
)

// Basic struct describes basic command type.
type Basic struct {
	// Trigger is a sequence of symbols message should start with to match with the command.
	trigger string
	// RootOnly determines if the command is supposed to be used by root only.
	RootOnly bool
	// PermissionsRequired is a slice of all permissions required by the command (but not subcommands).
	PermissionsRequired []int
	// Response is a string that will be sent to the user in response to the command.
	Response string
	// Description should contain short command description.
	Description string
	// Usage contains an example of the command usage.
	Usage string
}

func (c *Basic) Trigger() (trigger string) {
	return c.trigger
}

func (c *Basic) SetTrigger(trigger string) {
	c.trigger = trigger
}

func (c *Basic) HelpEmbed(sg *sugo.Instance, m *discordgo.Message) (embed *discordgo.MessageEmbed) {
	if c.Trigger() == "" || c.Description == "" || c.Usage == "" {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: "Developer of this command did not supply it with description. :frowning:",
			Color:       sugo.COLOR_WARNING,
		}
		return
	} else {
		embed = &discordgo.MessageEmbed{
			Title:       m.Content,
			Description: c.Description,
			Color:       sugo.COLOR_INFO,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Usage:",
					Value: fmt.Sprintf("`@%s` %s", sg.Self.Username, c.Usage),
				},
			},
		}
		return
	}

}

// Match checks if command signature matches the message content.
func (c *Basic) Match(sg *sugo.Instance, m *discordgo.Message) (matched bool, err error) {
	// By default command is not matched.
	matched = false
	if c.Trigger() == "" && m.Content == "" {
		return true, nil
	}

	if c.Trigger() != "" {
		if strings.HasPrefix(m.Content, c.Trigger()) {
			matched = true
			return true, nil
		}
	}
	return
}

// CheckPermissions checks message author and bot permissions if they match command's required permissions.
func (c *Basic) CheckPermissions(sg *sugo.Instance, m *discordgo.Message) (passed bool, err error) {
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

// Execute performs commands actions. For basic command it's just a simple text response.
func (c *Basic) Execute(sg *sugo.Instance, m *discordgo.Message) (err error) {
	_, err = c.RespondWithMention(sg, m, c.Response)
	if err != nil {
		return
	}
	return
}

// Responds to the channel without mention of the original message author.
func (c *Basic) Respond(sg *sugo.Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	message, err = sg.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		return
	}
	return
}

// Responds to the channel with the original message author mention.
func (c *Basic) RespondWithMention(sg *sugo.Instance, m *discordgo.Message, text string) (message *discordgo.Message, err error) {
	response_text := fmt.Sprintf("%s %s", helpers.UserAsMention(m.Author), text)
	message, err = sg.ChannelMessageSend(m.ChannelID, response_text)
	if err != nil {
		return
	}
	return
}
