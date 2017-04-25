# Sugo

Sugo is essentially a bot framework built on top of [discordgo bindings](https://github.com/bwmarrin/discordgo). There are few convenience features, such as:

- Helper functions (for example, to display discordgo.User or discordgo.Channel as mention for easier embedding in the messages)
- Ability to easily set up separate command repositories (such as commands for Elite Dangerous, commands for moderation etc) and activate only those you need.
- and eventually more... hopefully.

There are few repositories with the commands available:
- [sugo_commands_std](https://github.com/diraven/sugo_commands_std)
- well, for now there is only one. =)

The project is on rather early stages of development, so until first release there will be little effort to keep commands API backwards-compatible. Please, fork or checkout a commit that works for you.

If you have something to share, discuss or propose - please, do not hesitate to do so.

Typical bot initialization would look like the following:

```go
package main

import (
	"os"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo_commands_std"
)

func main() {
	// You can build your own simple custom command right here if you so desire.
	testCommand := sugo.Command{}
	testCommand.SetTrigger("test")
	testCommand.AddRequiredPermission(sugo.PermissionNone)
	testCommand.SetResponse("Passed!")
	sugo.Bot.RegisterCommand(&testCommand)

	// If you don't like default command trigger (for example if it clashes with some other one), you can change it like
	// so:
	sugo_commands_std.Info.SetTrigger("info") // Change "info" to whatever you see appropriate.
	sugo.Bot.RegisterCommand(&sugo_commands_std.Info)

	// And add some other commands to your bot.
	sugo.Bot.RegisterCommand(&sugo_commands_std.Greet)
	sugo.Bot.RegisterCommand(&sugo_commands_std.Help)
	sugo.Bot.RegisterCommand(&sugo_commands_std.Sys)
	sugo.Bot.RegisterCommand(&sugo_commands_std.Data)

	// Now just start the bot up and see what happens.
	// Make sure to provide at least token via SUGO_TOKEN environment variable.
	err := sugo.Bot.Startup(os.Getenv("SUGO_TOKEN"), os.Getenv("SUGO_ROOT_UID"))
	if err != nil {
		// TODO: Report error.
	}
}
```
