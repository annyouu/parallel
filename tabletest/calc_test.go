package calc

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string 
		a, b int
		expected int
	}{
		{"1+1", 1, 1, 2},
        {"2+3", 2, 3, 5},
        {"-1+1", -1, 1, 0},
        {"0+0", 0, 0, 0},
	}

	for _, tc := range tests {
		tc := tc // range変数の閉じ込め
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // サブテストを並列実行
			got := Add(tc.a, tc.b)
			if got != tc.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}