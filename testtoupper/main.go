package main

import (
	"strings"
	"testing"
)

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name string
		input string
		want string
	}{
		{"lowercase", "hello", "HELLO"},
		{"mixedcase", "HeLLo", "HELLO"},
		{"uppercase", "HELLO", "HELLO"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // 並列実行を宣言

			got := ToUpper(tc.input)
			if got != tc.want {
				t.Errorf("ToUpper(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}