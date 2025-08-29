package richerror

import "errors"

type Kind int

const (
	KindBadRequest Kind = iota + 1
	KindTooManyRequests
	KindGone
	KindInvalid
	KindUnauthorized
	KindUnexpected
)

type Operation string

type RichError struct {
	operation Operation
	wrapper   error
	message   string
	kind      Kind
}

func New(operation Operation) RichError {
	return RichError{operation: operation}
}

func (r RichError) WithWrapper(wrapper error) RichError {
	r.wrapper = wrapper
	return r
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind
	return r
}

func (r RichError) Error() string {
	return r.message
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	var re RichError
	ok := errors.As(r.wrapper, &re)
	if !ok {
		return r.wrapper.Error()
	}

	return re.Message()
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	var re RichError
	ok := errors.As(r.wrapper, &re)
	if !ok {
		return 0
	}

	return re.Kind()
}
