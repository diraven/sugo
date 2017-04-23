# sugo

Sugo is essentially a bot framework built on top of [discordgo bindings](https://github.com/bwmarrin/discordgo). There are few convenience features, such as:

- Helper functions (for example, to display discordgo.User or discordgo.Channel as mention for easier embedding in the messages)
- Ability to easily set up separate command repositories (such as commands for Elite Dangerous, commands for moderation etc) and activate only those you need.
- and eventually more... hopefully.

There are few repositories with the commands available:
- [sugo-commands-std](https://github.com/diraven/sugo-commands-std)
- well, for now there is only one. =)

The project is on rather early stages of development, so untill first release there will be little effort to keep commands API backwards-compatible. Please, fork or checkout a commit that works for you.

If you have something to share, discuss or propose - please, do not hesitate to do so.

Typical bot initialization would look like the following:

```go
package main

import (
	"os"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo-commands-std"
)

func main() {
	sugo.Bot.RegisterCommand(sugo_commands_std.Info)
	sugo.Bot.RegisterCommand(sugo_commands_std.Shutdown)

	err := sugo.Bot.Startup(os.Getenv("SUGO_TOKEN"), os.Getenv("SUGO_ROOT_UID"))
	if err != nil {
		// TODO: Report error.
	}
}
```
