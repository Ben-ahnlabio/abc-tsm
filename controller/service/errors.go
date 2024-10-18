package service

import "errors"

type ErrorString string

const (
	INVALID_INPUT string = "INVALID_INPUT"
)

var (
	ErrEncoding = errors.New(INVALID_INPUT)
)

type SvcErr struct {
	Text string
	Msg  string
}

func (e SvcErr) Error() string {
	return e.Msg
}

func InvalidInputError(err error) *SvcErr {
	return &SvcErr{
		Text: INVALID_INPUT,
		Msg:  err.Error(),
	}
}
