package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestB64CreateCommand(t *testing.T) {
	testConfigs := map[string]B64Config{
		"decode": {
			Encode: false,
			Text:   "dGhpc2lzZm9yYmFzZTY0dGVzdA",
		},
		"decodeErr": {
			Encode: false,
			Text:   "dGhpc2lzZm9yFzZTY0dGVzdA",
		},
		"encode": {
			Encode: true,
			Text:   "thisisforbase64test",
		},
	}
	for name, test := range testConfigs {
		t.Run(name, func(t *testing.T) {
			b64Cmd := B64CreateCommand(&test)
			assert.Equal(t, b64Cmd.Use, "b64", "Could not execute B64CreateCommand")
		})
	}
}

func TestB64Decode(t *testing.T) {
	text := "dGhpc2lzZm9yYmFzZTY0dGVzdA=="
	res, err := B64Decode(text)
	assert.Equal(t, err, nil, "Could not execute B64Decode()")
	assert.Equal(t, res, "thisisforbase64test", "Wrong output")
}

func TestB64Encode(t *testing.T) {
	text := "thisisforbase64test"
	res, err := B64Encode(text)
	assert.Equal(t, err, nil, "Could not execute B64Encode()")
	assert.Equal(t, res, "dGhpc2lzZm9yYmFzZTY0dGVzdA==", "Wrong output")
}
