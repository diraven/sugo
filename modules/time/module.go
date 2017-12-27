package time

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"time"
	"strings"
	"errors"
)

var defaultFormat = "2006.01.02 15:04  MST"

var timezones = tTimezones{}

func getLoc(sg *sugo.Instance, m *discordgo.Message, q string) (*time.Location, error) {
	// If query not specified - we fetch user/guild/default timezone.
	if q == "" {
		var err error
		if q, err = timezones.get(sg, m); err != nil {
			return nil, err
		}
	}

	// Try to find timezone in our offsets list.
	if offset, ok := offsets[strings.ToUpper(q)]; ok {
		return time.FixedZone(strings.ToUpper(q), offset), nil
	}

	// Try to search all other timezones in system tzdata.
	loc, err := time.LoadLocation(q)
	if err != nil {
		return nil, errors.New("timezone not found")
	}

	return loc, nil
}

// Module allows for time conversion operations.
var Module = &sugo.Module{
	RootCommand: rootCommand,
	Startup:     startup,
}
