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
		"happy path IPV6": {
			input:  "1::",
			expect: "[1::]:80",
		},
		"happy path IPV6:Port": {
			input:  "[1::]:22",
			expect: "[1::]:22",
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
		"happy path IPv6 for example.com": {
			input: "2606:2800:220:1:248:1893:25c8:1946",
		},
		"happy path IPv6 for google dns": {
			input: "[2001:4860:4860::8888]:53", // 8.8.8.8
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
