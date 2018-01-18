package time

import (
	"errors"
	"github.com/diraven/sugo"
	"strings"
	"time"
)

var defaultFormat = "2006.01.02 15:04  MST"

var timezones = tTimezones{}

func getLoc(sg *sugo.Instance, req *sugo.Request) (*time.Location, error) {
	// If query not specified - we fetch user/guild/default timezone.
	if req.Query == "" {
		var err error
		if req.Query, err = timezones.get(sg, req); err != nil {
			return nil, err
		}
	}

	// Try to find timezone in our offsets list.
	if offset, ok := offsets[strings.ToUpper(req.Query)]; ok {
		return time.FixedZone(strings.ToUpper(req.Query), offset), nil
	}

	// Try to search all other timezones in system tzdata.
	loc, err := time.LoadLocation(req.Query)
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
