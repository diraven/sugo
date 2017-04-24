package sugo

type Error struct {
	Text string
}

func (se Error) Error() (string) {
	return se.Text
}