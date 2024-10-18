package tsmutils

import "errors"

type ErrorString string

const (
	DECODING_ERROR string = "DECODING_ERROR"
)

var (
	ErrEncoding = errors.New(DECODING_ERROR)
)

type TsmUtilsErr struct {
	Text string
	Msg  string
}

func (e TsmUtilsErr) Error() string {
	return e.Msg
}

func DecodingError(err error) *TsmUtilsErr {
	return &TsmUtilsErr{
		Text: DECODING_ERROR,
		Msg:  err.Error(),
	}
}
