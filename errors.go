package zeroentd

import "fmt"

func NewErrConnection(e error) ErrConnection {
	return &connError{
		err: e,
	}
}

type ErrConnection interface {
	error
	IsConnError() bool
}

type connError struct {
	err error
}

func (impl *connError) IsConnError() bool {
	return true
}

func (impl *connError) Unwrap() error {
	return impl.err
}

func (impl *connError) Error() string {
	return "fluentd connection error"
}

type ErrJsonDecode interface {
	error
	IsJsonDecodeError() bool
}

func NewErrJsonDecode(e error) ErrJsonDecode {
	return &jsonDecodeError{
		err: e,
	}
}

type jsonDecodeError struct {
	err error
}

func (impl *jsonDecodeError) IsJsonDecodeError() bool {
	return true
}

func (impl *jsonDecodeError) Unwrap() error {
	return impl.err
}

func (impl *jsonDecodeError) Error() string {
	return fmt.Sprintf("failed json decode: %s", impl.err.Error())
}
