![code style](https://goreportcard.com/badge/github.com/diraven/sugo)

# Sugo

Sugo is a command-based bot framework built on top of [discordgo](https://github.com/bwmarrin/discordgo).

The idea of the project is to make simple wrapper around [discordgo](https://github.com/bwmarrin/discordgo) that allows fast and easy building of the commands-based bots without too much of a boilerplate code.

The project is on rather early stages of development, so until first release there will be little effort to keep commands API backwards-compatible. Please, fork or checkout a commit that works for you.

If you are worried about API being broken, please use one of the [released versions](https://github.com/diraven/sugo/releases) instead of regular `go get github.com/diraven/sugo`.

If you have something to share, discuss or propose - please, do not hesitate to do so via [issues](https://github.com/diraven/sugo/issues).

## Quickstart

The simpler command would look like the following:

```go
package main

import (
	"github.com/diraven/sugo"
)

// Make a command.
var cmd = &sugo.Command{
	Trigger:     "ping",
	Description: "just responds with pong!",
	Execute: func(req *sugo.Request) error {
		_, err := req.RespondInfo("", "pong!")
		return err
	},
}

func main() {
	// Create new bot instance.
	bot := sugo.New()

	// Set bot trigger (not required, by default the bot will react to the messages starting with the bot @mention).
	bot.Trigger = "."

	// Initialize modules.
	bot.AddCommand(cmd)

	// Start the bot.
	if err := bot.Startup("bot YOURTOKENSTRING"); err != nil {
		bot.HandleError(err)
	}
}
```

## Advanced Usage

### Underlying *discordgo.Session access

You may need it to perform numerous API interactions such as assigning user roles etc. Inside commands it is available via Sugo instance `req.Sugo.Session`, see [godoc](https://godoc.org/github.com/diraven/sugo#Instance) for details. 

### Custom embeds

If you need some custom message type and regular `req.Respond*` is not enough, you can send messages directly using lower level `*discordgo.Session`, like this:
```go
package main

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

// Make a command.
var cmd = &sugo.Command{
	Trigger:     "ping",
	Description: "just responds with pong!",
	Execute: func(req *sugo.Request) error {
		embed := &discordgo.MessageEmbed{
			Title: "Command Response",
			Description: "Pong!",
			Color: sugo.ColorDefault,
		}
		
		_, err := req.Sugo.Session.ChannelMessageSendEmbed(req.Channel.ID, embed)
		return err
	},
}

func main() {
	// Create new bot instance.
	bot := sugo.New()

	// Add command.
	bot.AddCommand(cmd)

	// Start the bot.
	if err := bot.Startup("bot YOURTOKENSTRING"); err != nil {
		bot.HandleError(err)
	}
}
```
 
 See [discordgo documentation](https://godoc.org/github.com/bwmarrin/discordgo) for additional details.

### Permissions

Bot comes with a simplistic permissions system. It utilizes discord permissions. Command can be restricted to the users that have specified permissions.

For example if you want the command to be restricted to users that can ban members AND add reactions, command declaration will look like this:

```go
package main

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
)

// Make a command.
var cmd = &sugo.Command{
	Trigger:     "ping",
	Description: "just responds with pong!",
	PermissionsRequired: discordgo.PermissionBanMembers | discordgo.PermissionAddReactions,
	Execute: func(req *sugo.Request) error {
		_, err := req.RespondInfo("", "pong!")
		return err
	},
}

func main() {
	// Create new bot instance.
	bot := sugo.New()

	// Add command.
	bot.AddCommand(cmd)

	// Start the bot.
	if err := bot.Startup("bot YOURTOKENSTRING"); err != nil {
		bot.HandleError(err)
	}
}
```

### More info

See [godoc](https://godoc.org/github.com/diraven/sugo) and [command modules examples](https://github.com/diraven/sugo/tree/master/examples).

Typical bot initialization would look like the following:

```go
package main

import (
	"github.com/diraven/sugo"
	"github.com/diraven/sugo/examples/test"
	"github.com/diraven/sugo/examples/help"
	"github.com/diraven/sugo/examples/info"
)

func main() {
	// Create new bot instance.
	bot := sugo.New()

	// Set bot trigger.
	bot.Trigger = "."

	// Initialize modules.
	test.Init(bot)
	help.Init(bot)
	info.Init(bot)

	// Start the bot.
	if err := bot.Startup("bot YOURTOKENSTRING"); err != nil {
		bot.HandleError(err)
	}
}
```

Compile and run your bot, then try to write something like `.test responses` in the channel bot has both read and write access to.