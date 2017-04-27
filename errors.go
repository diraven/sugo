package sugo

// Error represents any Sugo error...
type Error struct {
	Text string
}

func (se Error) Error() string {
	return se.Text
}
