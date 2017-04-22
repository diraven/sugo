package errors

type SugoError struct {
	Text string
}

func (se SugoError) Error() (string) {
	return se.Text
}