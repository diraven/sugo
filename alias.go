package sugo

type alias struct {
	from string
	to   string
}

type aliases []alias

// AddAlias adds function that will be called on bot startup.
func (sg *Instance) AddAlias(from, to string) {
	*sg.aliases = append(*sg.aliases, alias{from: from, to: to})
}
