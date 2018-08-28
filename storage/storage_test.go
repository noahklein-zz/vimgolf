package storage

import (
	"testing"
)

func Test_isValidUsername(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "too short", input: "Jo", want: false},
		{name: "invalid chars", input: "~josh", want: false},
		{name: "valid", input: "__Jeff47", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUsername(tt.input); got != tt.want {
				t.Errorf("isValidUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
