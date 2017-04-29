package sugo

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo/helpers"
	"log"
	"strings"
)

// Command struct describes basic command type.
type Command struct {
	// trigger is a sequence of symbols message should start with to match with the command.
	Trigger string
	// rootOnly determines if the command is supposed to be used by root only.
	RootOnly bool
	// permissionsRequired is a slice of all permissions required by the command (but not subcommands).
	Permissions []int
	// response is a string that will be sent to the user in response to the command.
	TextResponse string
	// embedResponse is a *discordgo.MessageEmbed, if set - has priority over text response.
	EmbedResponse *discordgo.MessageEmbed
	// description should contain short command description.
	Description string
	// usage contains an example of the command usage.
	Usage string
	// subCommands contains all subcommands of the given command.
	SubCommands []*Command

	// parentCommand contains command, which is parent for this one
	parent *Command
	// subCommandsTriggers contains all triggers of subcommands for the help to refer to.
	subCommandsTriggers []string

	Execute   func(c *Command, sg *Instance, m *discordgo.Message) (err error)
	HelpEmbed func(c *Command, sg *Instance) (embed *discordgo.MessageEmbed, err error)
	Startup   func(c *Command, sg *Instance) (err error)
	Teardown  func(c *Command, sg *Instance) (err error)
}

// cmdStartup is internal function called for each command on bot startup.
func cmdStartup(c *Command, sg *Instance) (err error) {
	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Build command triggers cache.
		if v.Trigger != "" {
			c.subCommandsTriggers = append(c.subCommandsTriggers, v.Trigger)
		}

		// Check if command is already registered elsewhere.
		if v.parent != nil {
			return Error{
				fmt.Sprintf("The subcommand is already registered elsewhere: %s", cmdFullUsage(c, sg)),
			}
		}
		// Set command parent.
		v.parent = c

		// Run system startup for subcommand.
		cmdStartup(v, sg)
	}

	// Run public startup for command if set.
	if c.Startup != nil {
		err = c.Startup(c, sg)
		if err != nil {
			return
		}
	}

	return
}

// cmdTeardown is internal function called for each command on bot graceful shutdown.
func cmdTeardown(c *Command, sg *Instance) error {
	var err error

	// For every subcommand (if any):
	for _, v := range c.SubCommands {
		// Here be some internal code to tear commands down... some day. May be.

		// Run system startup for subcommand.
		err = cmdTeardown(v, sg)
		if err != nil {
			log.Fatal("Command teardown error: ", err)
		}
	}

	// Run public teardown for command if set.
	if c.Teardown != nil {
		err = c.Teardown(c, sg)
		if err != nil {
			return Error{fmt.Sprintf("Command custom teardown error: %s\n", err)}
		}
	}
	return nil
}

// cmdPath returns all the triggers from parent commands from outermost to innermost parent.
func cmdPath(c *Command) (value string) {
	if c.parent != nil {
		return cmdPath(c.parent) + " " + c.Trigger
	}
	return c.Trigger
}

// cmdFullUsage returns full command usage including all parent triggers.
func cmdFullUsage(c *Command, sg *Instance) (value string) {
	return helpers.UserAsMention(sg.Self) + cmdPath(c) + " " + c.Usage
}

// cmdHelpEmbed is a default implementation of help embed builder.
func cmdHelpEmbed(c *Command, sg *Instance) (embed *discordgo.MessageEmbed) {
	if c.Trigger == "" || c.Description == "" {
		embed = &discordgo.MessageEmbed{
			Title:       cmdPath(c),
			Description: "Developer of this command did not supply it with description. :frowning:",
			Color:       ColorWarning,
		}
		return embed
	}
	embed = &discordgo.MessageEmbed{
		Title:       cmdPath(c),
		Description: c.Description,
		Color:       ColorInfo,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Usage:",
				Value: cmdFullUsage(c, sg),
			},
		},
	}
	if len(c.SubCommands) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Subcommands:",
				Value: strings.Join(c.subCommandsTriggers, ", "),
			}, &discordgo.MessageEmbedField{
				Name:  "To get help on 'subcommand' type:",
				Value: fmt.Sprintf("`@%s` help %s subcommand", sg.Self.Username, c.Trigger),
			})
	}
	return embed

}

// cmdMatch is a default matching function that checks if command trigger matches the start of message content.
func cmdMatch(c *Command, sg *Instance, m *discordgo.Message) (matched bool, err error) {
	// By default command is not matched.
	matched = false

	// If trigger is not set, check if command is empty.
	if c.Trigger == "" && m.Content == "" {
		return true, nil
	}

	// Trigger is set, see if it's in the message.
	if c.Trigger != "" {
		if strings.HasPrefix(m.Content, c.Trigger) {
			matched = true
			return true, nil
		}
	}
	return
}

// checkCheckPermissions checks message author and bot permissions if they match command's required permissions.
func cmdCheckPermissions(c *Command, sg *Instance, m *discordgo.Message) (passed bool, err error) {
	// By default command is not allowed.
	passed = false

	// For security reasons - every command should have at least one permission set explicitly.
	if len(c.Permissions) == 0 {
		err = Error{Text: "Command has no Permissions[]!"}
		return
	}

	// Get channel to check for permissions.
	channel, err := sg.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// Calculate compound permission.
	var compoundPerm int
	for _, perm := range c.Permissions {
		compoundPerm |= perm
	}

	// Make sure bot has the permission required.
	botHasPerm, err := sg.BotHasPermission(compoundPerm, channel)
	if err != nil {
		return
	}
	if !(botHasPerm) {
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
	userHasPerm, err := sg.UserHasPermission(compoundPerm, m.Author, channel)
	if err != nil {
		return
	}
	if !(userHasPerm) {
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

// cmdExecute is a default command execution function.
func cmdExecute(c *Command, sg *Instance, m *discordgo.Message) (err error) {
	var actionPerformed bool

	if c.Execute != nil {
		// Run custom command Execute if set.
		err = c.Execute(c, sg, m)
		if err != nil {
			return
		}
		actionPerformed = true
	}

	if c.TextResponse != "" {
		// Send command text response if set.
		_, err = sg.RespondTextMention(m, c.TextResponse)
		if err != nil {
			return
		}
		actionPerformed = true
	}

	if c.EmbedResponse != nil {
		// Send command embed response if set.
		_, err = sg.RespondEmbed(m, c.EmbedResponse)
		if err != nil {
			return
		}
		actionPerformed = true
	}

	if !actionPerformed {
		if len(c.SubCommands) > 0 {
			// If there is at least one subcommand and no other actions taken - explain it to the user.
			_, err = sg.RespondTextMention(m, "This command itself does not seem to do anything. Try subcommands instead.")
			if err != nil {
				return
			}
		}

		// We did nothing and there are no subcommands...
		_, err = Bot.RespondTextMention(m, "Looks like this command just does nothing... What is it here for?")
		if err != nil {
			return
		}
	}

	return
}
