package wrap_test

import (
	"fmt"
	"testing"

	"github.com/progfay/go-wrap/wrap"
)

func Test_Pipe(t *testing.T) {
	t.Run("all piped functions must be called", func(t *testing.T) {
		called1, called2, called3 := false, false, false
		fn1 := func() {
			called1 = true
		}
		fn2 := func() {
			called2 = true
		}
		fn3 := func() {
			called3 = true
		}

		var pipedFn func()
		err := wrap.Pipe(&pipedFn, fn1, fn2, fn3)
		if err != nil {
			t.Error(err)
		}

		pipedFn()
		if !called1 || !called2 || !called3 {
			t.Error("All piped functions must be called")
		}
	})

	t.Run("process values in order by functions", func(t *testing.T) {
		fn1 := func(a, b int) (int, int) {
			return a + b, a * b
		}
		fn2 := func(a, b int) (string, string) {
			return fmt.Sprintf("%v", a), fmt.Sprintf("%v", b)
		}
		fn3 := func(a, b string) string {
			return fmt.Sprintf("%v,%v", a, b)
		}

		var pipedFn func(a, b int) string
		err := wrap.Pipe(&pipedFn, fn1, fn2, fn3)
		if err != nil {
			t.Error(err)
		}

		ret := pipedFn(2, 3)
		if ret != "5,6" {
			t.Errorf("Returned value is incorrect, want: %v, got: %v", "5,6", ret)
		}
	})

	t.Run("first argument must be function pointer", func(t *testing.T) {
		fn1 := func() {}
		fn2 := func() {}
		fn3 := func() {}
		err := wrap.Pipe(fn1, fn2, fn3)
		if err == nil {
			t.Error("cannot set other than function pointer for first argument")
		}

		intValue := 0
		fn4 := func() {}
		fn5 := func() {}
		err = wrap.Pipe(&intValue, fn4, fn5)
		if err == nil {
			t.Error("cannot set other than function pointer for first argument")
		}
	})

	t.Run("all the arguments after the second must be function", func(t *testing.T) {
		var dst func()
		err := wrap.Pipe(&dst, 1)
		if err == nil {
			t.Error("all the arguments after the second must be function")
		}

		fn := func() {}
		err = wrap.Pipe(&dst, fn, 1)
		if err == nil {
			t.Error("all the arguments after the second must be function")
		}
	})

	t.Run("all the arguments after the second must have returned values type same as arguments type of next function", func(t *testing.T) {
		var dst func(a int64)
		fn1 := func(a int64) {}
		fn2 := func(a int64) {}
		err := wrap.Pipe(&dst, fn1, fn2)
		if err == nil {
			t.Error("all the arguments after the second must have returned values type same as arguments type of next function")
		}
	})

	t.Run("second argument must have arguments type same as arguments type of value of first argument", func(t *testing.T) {
		var dst func(a int32) int64
		fn1 := func(a int64) int64 { return 0 }
		err := wrap.Pipe(&dst, fn1)
		if err == nil {
			t.Error("second argument must have arguments type same as arguments type of value of first argument")
		}

		fn2 := func(a int64) int64 { return 0 }
		err = wrap.Pipe(&dst, fn1, fn2)
		if err == nil {
			t.Error("second argument must have arguments type same as arguments type of value of first argument")
		}
	})

	t.Run("last argument must have returned values type same as returned values type of value of first argument", func(t *testing.T) {
		var dst func(a int64) int32
		fn1 := func(a int64) int64 { return 0 }
		err := wrap.Pipe(&dst, fn1)
		if err == nil {
			t.Error("last argument must have returned values type same as returned values type of value of first argument")
		}

		fn2 := func(a int64) int64 { return 0 }
		err = wrap.Pipe(&dst, fn1, fn2)
		if err == nil {
			t.Error("last argument must have returned values type same as returned values type of value of first argument")
		}
	})
}
