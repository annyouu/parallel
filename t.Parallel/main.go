package parallel_test

import (
	"fmt"
	"testing"
	"time"
)

func trace(name string) func() {
	fmt.Printf("%s entered\n", name)
	return func() {
		fmt.Printf("%s returned\n", name)
	}
}

func Test_Main1(t *testing.T) {
	defer trace("Test_Main1")()

	t.Run("Sub1", func(t *testing.T) {
		defer trace("Sub1")()
		t.Parallel()
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Sub2", func(t *testing.T) {
		defer trace("Sub2")()
		t.Parallel()
		time.Sleep(100 * time.Millisecond)
	})
}

func Test_Main2(t *testing.T) {
	defer trace("Test_Main2")()
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
}

func Test_Main3(t *testing.T) {
	defer trace("Test_Main3")()
	time.Sleep(100 * time.Millisecond)
}

