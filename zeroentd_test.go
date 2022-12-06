package zeroentd

import (
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
