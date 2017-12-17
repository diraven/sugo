package triggers

import (
	"github.com/diraven/sugo"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func onBeforeBotTriggerDetect(sg *sugo.Instance, m *discordgo.Message, q string) (string, error) {
	var err error

	g, err := sg.GuildFromMessage(m)
	if err != nil {
		return q, err
	}

	q = strings.TrimSpace(q)

	trigger := triggers.get(sg, g.ID)

	// If trigger for the given guild is set and present in the query:
	if trigger != "" && strings.HasPrefix(q, trigger) {
		// Replace trigger with bot mention for it to be detected as bot trigger.
		q = strings.Replace(q, trigger, sg.Self.Mention(), 1)
	}

	// Return our resulting query string.
	return q, nil
}
