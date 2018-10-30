package sugo

// Error is a public error that will be sent to the user.
type Error struct {
	request *Request
	text    string
}

func (e *Error) Error() string {
	return e.text
}

func WrapError(req *Request, cause error) (err error) {
	if cause != nil {
		err = &Error{request: req, text: err.Error()}
	}
	return
}

func NewError(req *Request, text string) (err error) {
	return &Error{request: req, text: text}
}

func NewNotImplementedError(req *Request) (err error) {
	return &Error{request: req, text: "this functionality is not implemented yet"}
}

func NewBadCommandUsageError(req *Request) (err error) {
	return &Error{request: req, text: "command used incorrectly, check `.help` for details"}
}

func NewGuildOnlyError(req *Request) (err error) {
	return &Error{request: req, text: "this command can only be used in a guild channel"}
}

func NewCommandNotFoundError(req *Request) (err error) {
	return &Error{request: req, text: "command not found"}
}
