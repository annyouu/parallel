package target

import (
	"testing"
)

func TestParent(t *testing.T) {
	t.Run("Subtest1", func(t *testing.T) {
		t.Parallel() // サブテストでParallelを呼び出す
	})
	t.Run("Subtest2", func(t *testing.T) {
		// サブテストでParallelを呼ばない
	})
}

func SubTestParent(t *testing.T) {
	t.Run("SubtestA", func(t *testing.T) {
		// サブテストでParallelは呼ばれない
	})
}

func TestFoo(t *testing.T) {
    t.Run("sub", func(t *testing.T) {
        t.Parallel()
    })
}