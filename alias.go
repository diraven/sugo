package sugo

type alias struct {
	from string
	to   string
}

type aliases []alias

// AddAlias adds command alias. If command was not found - the query string will be searched for "from" string, and if
// found "from" will be replaced with "to", then another command search will be performed. Only one replacement
// can happen for one request.
func (sg *Instance) AddAlias(from, to string) {
	*sg.aliases = append(*sg.aliases, alias{from: from, to: to})
}
