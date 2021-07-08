package reachable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReachableChecker_resolve(t *testing.T) {
	testcases := map[string]struct {
		input   string
		expect  string
		wantErr bool
	}{
		"happy path URL": {
			input:  "https://example.com",
			expect: "example.com:80",
		},
		"happy path websocket address": {
			input:  "wss://example.com:8443",
			expect: "example.com:8443",
		},
		"happy path Domain:Port": {
			input:  "example.com:8080",
			expect: "example.com:8080",
		},
		"happy path IP:Port": {
			input:  "127.0.0.1:22",
			expect: "127.0.0.1:22",
		},
		"sad path empty": {
			input:   "",
			wantErr: true,
		},
		"sad path invalid URL": {
			input:   "https://example. com",
			wantErr: true,
		},
	}

	checker := New(nil)

	for name, tt := range testcases {
		t.Run(name, func(t *testing.T) {
			actual, err := checker.resolve(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestReachableChecker_Check(t *testing.T) {
	testcases := map[string]struct {
		input   string
		wantErr bool
	}{
		"happy path URL": {
			input: "https://example.com",
		},
	}
	for name, tt := range testcases {
		t.Run(name, func(t *testing.T) {
			checker := New(nil)
			err := checker.Check(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
