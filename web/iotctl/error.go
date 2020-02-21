package iotctl

type Error struct {
	Err   error
	Stack []byte
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) StackTrace() []byte {
	return e.Stack
}

type StackTrace interface {
	StackTrace() []byte
}
