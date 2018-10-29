package sugo

// Error is a public error that will be sent to the user.
type Error struct {
	request *Request
	text    string
}

func (e *Error) Error() string {
	return e.text
}

func WrapError(req *Request, err error) *Error {
	return &Error{request: req, text: err.Error()}
}

func NewError(req *Request, text string) *Error {
	return &Error{request: req, text: text}
}

func NotImplementedError(req *Request) *Error {
	return &Error{request: req, text: "this functionality is not implemented yet"}
}

func NewBadCommandUsageError(req *Request) *Error {
	return &Error{request: req, text: "command used incorrectly, check `.help` for more details"}
}

func NewGuildOnlyError(req *Request) *Error {
	return &Error{request: req, text: "this command can only be used in a guild channel"}
}

func NewCommandNotFoundError(req *Request) *Error {
	return &Error{request: req, text: "command not found"}
}
