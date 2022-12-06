package zeroentd

import (
	"bytes"
	"encoding/json"
	"unsafe"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/pkg/errors"
)

type fluentIface interface {
	Post(tag string, message interface{}) error
	Close() error
}

func New(tag string, cfg ...fluent.Config) (*Writer, error) {
	if len(cfg) == 0 {
		client, err := fluent.New(fluent.Config{})
		if err != nil {
			return nil, errors.WithMessage(NewErrConnection(err), "cannot connect fluentd")
		}
		return &Writer{
			client: client,
			tag:    tag,
		}, nil
	}
	client, err := fluent.New(cfg[0])
	if err != nil {
		return nil, errors.WithMessage(NewErrConnection(err), "cannot connect fluentd")
	}
	return &Writer{
		client: client,
		tag:    tag,
	}, nil
}

type Writer struct {
	client fluentIface
	tag    string
}

func (w *Writer) Close() error {
	err := w.client.Close()
	if err != nil {
		return errors.WithMessagef(err, "fluentd client could not be closed")
	}
	return nil
}

func (w *Writer) Write(data []byte) (int, error) {
	r := bytes.NewReader(data)
	packet := map[string]any{}
	if err := json.NewDecoder(r).Decode(&packet); err != nil {
		return 0, errors.WithMessage(NewErrJsonDecode(err), "failed to parse zerolog log string")
	}
	if err := w.client.Post(w.tag, packet); err != nil {
		return 0, errors.WithMessagef(err, "failed to send log: %s", bytesToStrUnsafe(data))
	}
	return len(data), nil
}

func bytesToStrUnsafe(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
