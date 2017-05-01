package sugo

// sError represents any Sugo error...
type sError struct {
	s string
}

func (e sError) Error() string {
	return e.s
}