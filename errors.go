package zeroentd

func NewErrConnection(e error) *ErrConnection {
	return &ErrConnection{
		err: e,
	}
}

type ErrConnection struct {
	err error
}

func (impl *ErrConnection) Unwrap() error {
	return impl.err
}

func (impl *ErrConnection) Error() string {
	return "fluentd connection error"
}
