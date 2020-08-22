package wrap_test

import (
	"testing"

	"github.com/progfay/go-wrap/wrap"
)

func Test_Before(t *testing.T) {
	t.Run("wrap function with no args and return values", func(t *testing.T) {
		beforeFnIsCalled, fnIsCalled := false, false
		beforeFn := func() {
			beforeFnIsCalled = true
		}
		fn := func() {
			if !beforeFnIsCalled {
				t.Error("Third argument must be called before second argument")
			}
			fnIsCalled = true
		}

		var wrappedFn func()
		err := wrap.Before(&wrappedFn, fn, beforeFn)
		if err != nil {
			t.Error(err)
		}

		wrappedFn()
		if !fnIsCalled {
			t.Error("Second argument must be called")
		}
	})

	t.Run("wrap function with many args and return values", func(t *testing.T) {
		beforeFnIsCalled, fnIsCalled := false, false
		beforeFn := func() {
			beforeFnIsCalled = true
		}
		fn := func(a, b int) (int, int) {
			if !beforeFnIsCalled {
				t.Error("Third argument must be called before second argument")
			}
			fnIsCalled = true
			return a + b, a * b
		}

		var wrappedFn func(a, b int) (int, int)
		err := wrap.Before(&wrappedFn, fn, beforeFn)
		if err != nil {
			t.Error(err)
		}

		ret1, ret2 := wrappedFn(2, 3)
		if !fnIsCalled {
			t.Error("Second argument must be called")
		}

		if ret1 != 5 || ret2 != 6 {
			t.Error("returned value is incorrect")
		}
	})

	t.Run("first argument must be function pointer", func(t *testing.T) {
		fn1 := func() {}
		fn2 := func() {}
		fn3 := func() {}
		err := wrap.Before(fn1, fn2, fn3)
		if err == nil {
			t.Error("cannot set other than function pointer for first argument")
		}

		intValue := 0
		fn4 := func() {}
		fn5 := func() {}
		err = wrap.Before(&intValue, fn4, fn5)
		if err == nil {
			t.Error("cannot set other than function pointer for first argument")
		}
	})

	t.Run("second argument must be function", func(t *testing.T) {
		var dst func()
		fn := func() {}
		err := wrap.Before(dst, 1, fn)
		if err == nil {
			t.Error("cannot set other than function for second argument")
		}
	})

	t.Run("value of first argument must be same type as second argument", func(t *testing.T) {
		var dst func(a int32)
		fn := func(a int64) {}
		err := wrap.Before(&dst, fn, func() {})
		if err == nil {
			t.Error("value of first argument must be same type as second argument")
		}
	})
}
