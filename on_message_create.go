package sugo

import (
	"strings"
	"github.com/bwmarrin/discordgo"
	"context"
	"errors"
)

// onMessageCreate contains all the message processing logic for the bot.
func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	var err error                  // Used to capture and report errors.
	var ctx = context.Background() // Root context.
	var command *Command           // Used to store the command we will execute.
	var q = mc.Content             // Command query string.

	// Make sure we are in the correct bot instance.
	if Bot.Session != s {
		Bot.HandleError(errors.New("Bot session error:" + err.Error()))
		Bot.Shutdown()
	}

	// Make sure message author is not a bot.
	if mc.Author.Bot {
		return
	}

	// OnBeforeBotTriggerDetect entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeBotTriggerDetect != nil {
			q, err = module.OnBeforeBotTriggerDetect(Bot, mc.Message, q)
			if err != nil {
				Bot.HandleError(errors.New("OnBeforeMentionDetect error: " + err.Error() + " (" + q + ")"))
			}
		}
	}

	// If bot nick was changed on the server - it will have ! in it's mention, so we need to remove that in order
	// for mention detection to work right.
	if strings.HasPrefix(q, "<@!") {
		q = strings.Replace(q, "<@!", "<@", 1)
	}

	// Make sure message starts with bot mention.
	if strings.HasPrefix(strings.TrimSpace(q), Bot.Self.Mention()) {
		// Remove bot trigger from the string.
		q = strings.TrimSpace(strings.TrimPrefix(q, Bot.Self.Mention()))
	} else {
		return
	}

	// Fill context with necessary data.
	// Get Channel.
	channel, err := Bot.ChannelFromMessage(mc.Message)
	if err != nil {
		Bot.HandleError(err)
	}
	// Save into context.
	ctx = context.WithValue(ctx, CtxKey("channel"), channel)

	// Get Guild.
	guild, err := Bot.GuildFromMessage(mc.Message)
	if err != nil {
		Bot.HandleError(err)
	}
	// Save into context.
	ctx = context.WithValue(ctx, CtxKey("guild"), guild)

	// OnBeforeCommandSearch entry point for Modules.
	for _, module := range Bot.Modules {
		if module.OnBeforeCommandSearch != nil {
			q, err = module.OnBeforeCommandSearch(Bot, mc.Message, q)
			if err != nil {
				Bot.HandleError(errors.New("OnBeforeCommandSearch error: " + err.Error() + " (" + q + ")"))
			}
		}
	}

	// Search for applicable command.
	command, err = Bot.FindCommand(mc.Message, q)
	if err != nil {
		// Unhandled error in command.
		Bot.HandleError(errors.New("Bot command search error: " + err.Error() + " (" + q + ")"))
		Bot.Shutdown()
	}
	if command != nil {
		// Remove command trigger from message string.
		q = strings.TrimSpace(strings.TrimPrefix(q, command.Path()))

		// And execute command.
		err = command.execute(ctx, q, Bot, mc.Message)
		if err != nil {
			if strings.Contains(err.Error(), "\"code\": 50013") {
				// Insufficient permissions, bot configuration issue.
				Bot.HandleError(errors.New("Bot permissions error: " + err.Error() + " (" + q + ")"))
			} else {
				// Other discord errors.
				Bot.HandleError(errors.New("Bot command execute error: " + err.Error() + " (" + q + ")"))
				Bot.Shutdown()
			}
		}
		return
	}

	Bot.RespondCommandNotFound(mc.Message)

	// Command not found.
}
