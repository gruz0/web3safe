package app

type Error struct {
	message  string
	exitCode int
}

func NewError(message string, exitCode int) *Error {
	return &Error{
		message:  message,
		exitCode: exitCode,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) ExitCode() int {
	return e.exitCode
}
