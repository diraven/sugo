package triggers

import (
	"github.com/diraven/sugo"
	"strings"
)

func onBeforeBotTriggerDetect(sg *sugo.Instance, req *sugo.Request) error {
	// Trigger detection only works in guild text channels, so we are safe to assume we are in guild channel at this
	// point.

	// Cleanup query.
	req.Query = strings.TrimSpace(req.Query)

	// Get bot trigger.
	trigger := triggers.get(sg, req.Guild.ID)

	// If trigger for the given guild is set and present in the query:
	if trigger != "" && strings.HasPrefix(req.Query, trigger) {
		// Replace trigger with bot mention for it to be detected as bot trigger.
		req.Query = strings.Replace(req.Query, trigger, sg.Self.Mention(), 1)
	}

	return nil
}
