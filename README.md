![code style](https://goreportcard.com/badge/github.com/diraven/sugo)

# Sugo

Sugo is a command-based bot framework built on top of [discordgo](https://github.com/bwmarrin/discordgo).

The idea of the project is to make simple wrapper around [discordgo](https://github.com/bwmarrin/discordgo) that allows fast and easy building of the commands-based bots without too much of a boilerplate code.

The project is on rather early stages of development, so until first release there will be little effort to keep commands API backwards-compatible. Please, fork or checkout a commit that works for you.

If you are worried about API being broken, please use one of the [released versions](https://github.com/diraven/sugo/releases) instead of regular `go get github.com/diraven/sugo`.

If you have something to share, discuss or propose - please, do not hesitate to do so via [issues](https://github.com/diraven/sugo/issues).

## Usage

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

Compile and run your bot, then try to write something like `.test` in the channel bot has both read and write access to.