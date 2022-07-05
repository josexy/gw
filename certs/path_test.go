package certs

import (
	"testing"
)

func Test_Path(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"public.crt"},
		{"server.pem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(Path(tt.name))
		})
	}
}
