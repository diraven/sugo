![code style](https://goreportcard.com/badge/github.com/diraven/sugo)

# Sugo

Sugo is essentially a bot framework built on top of [discordgo bindings](https://github.com/bwmarrin/discordgo). There are few convenience features, such as:

- Helper functions (for example, to display discordgo.User or discordgo.Channel as mention for easier embedding in the messages)
- Ability to easily set up separate command repositories (such as commands for Elite Dangerous, commands for moderation etc) and activate only those you need.
- and eventually more... hopefully.

You can pull pull request your own commands or changes to the existing ones in the [sugo_contrib](https://github.com/diraven/sugo_contrib) repo.

The project is on rather early stages of development, so until first release there will be little effort to keep commands API backwards-compatible. Please, fork or checkout a commit that works for you.

If you have something to share, discuss or propose - please, do not hesitate to do so.

Typical bot initialization would look like the following:

```go
package main

import (
	"os"
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/modules/aliases"
	"github.com/diraven/sugo/modules/test"
	"github.com/diraven/sugo/modules/permissions"
	"github.com/diraven/sugo/modules/help"
	"github.com/diraven/sugo/modules/elite_dangerous"
	"github.com/diraven/sugo/modules/greet"
	"github.com/diraven/sugo/modules/guild_wars2"
	"github.com/diraven/sugo/modules/clean"
	"github.com/diraven/sugo/modules/info"
	"github.com/diraven/sugo/modules/public_roles"
	"github.com/diraven/sugo/modules/feeds"
	"github.com/diraven/sugo/modules/sys"
)

func main() {
	sugo.Bot.Modules = []*sugo.Module{
		aliases.Module,
		clean.Module,
		elite_dangerous.Module,
		feeds.Module,
		greet.Module,
		guild_wars2.Module,
		help.Module,
		info.Module,
		permissions.Module,
		public_roles.Module,
		sys.Module,
		test.Module,
	}

	if err := sugo.Bot.Startup(os.Getenv("SUGO_TOKEN"), os.Getenv("SUGO_ROOT_UID")); err != nil {
		// TODO: Report errors
	}
}
```
