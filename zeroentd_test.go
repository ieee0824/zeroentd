package zeroentd

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToStrUnsafe(t *testing.T) {
	tests := []struct {
		in string
	}{
		{
			in: "aaaa",
		},
		{
			in: "本日は晴天なり",
		},
		{
			in: "🍺",
		},
		{
			in: "",
		},
	}

	for _, test := range tests {
		r := bytesToStrUnsafe([]byte(test.in))
		assert.Equal(t, test.in, r)
	}
}

type dummyFluent struct {
	err error
}

func (impl *dummyFluent) Post(tag string, message interface{}) error {
	return impl.err
}

func (impl *dummyFluent) Close() error {
	return nil
}

func isErrType(e error, target error) bool {
	switch e.(type) {
	case ErrConnection:
		_, ok := e.(ErrConnection)
		return ok
	case ErrJsonDecode:
		_, ok := e.(ErrJsonDecode)
		return ok
	default:
		if e == target {
			return true
		}
		ure := errors.Unwrap(e)
		if ure == nil {
			return false
		}
		return isErrType(ure, target)
	}
}

var dummyError = errors.New("dummy error")

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name         string
		fluentClient *dummyFluent
		tag          string
		data         []byte
		isError      bool
		err          error
	}{
		{
			name:         "success",
			fluentClient: &dummyFluent{},
			tag:          "success.test.logger",
			data: func() []byte {
				bin, _ := json.Marshal(map[string]any{
					"foo": 1,
					"bar": "b",
					"baz": struct{}{},
				})
				return bin
			}(),
			isError: false,
			err:     nil,
		},
		{
			name:         "write not json",
			fluentClient: &dummyFluent{},
			tag:          "write.not.json.test.logger",
			data:         []byte("hoge"),
			isError:      true,
			err:          &jsonDecodeError{err: errors.New("not json")},
		},
		{
			name: "fluentd error",
			fluentClient: &dummyFluent{
				err: dummyError,
			},
			tag: "fluentd.err.test.logger",
			data: func() []byte {
				bin, _ := json.Marshal(map[string]any{
					"foo": 1,
					"bar": "b",
					"baz": struct{}{},
				})
				return bin
			}(),
			isError: true,
			err:     dummyError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &Writer{
				client: &dummyFluent{
					err: test.err,
				},
				tag: test.tag,
			}

			n, err := w.Write(test.data)
			if !test.isError {
				assert.Nil(t, err)
				assert.Equal(t, len(test.data), n)
			} else {
				if !assert.True(t, isErrType(err, test.err)) {
					t.Logf("error is %t", err)
				} else {
					assert.Equal(t, 0, n)
				}
			}
		})
	}
}
