package wrap_test

import (
	"testing"

	"github.com/progfay/go-wrap/wrap"
)

func Test_After(t *testing.T) {
	t.Run("wrap function with no args and return values", func(t *testing.T) {
		afterFnIsCalled, fnIsCalled := false, false
		fn := func() {
			fnIsCalled = true
		}
		afterFn := func() {
			if !fnIsCalled {
				t.Error("Second argument must be called before third argument")
			}
			afterFnIsCalled = true
		}

		var wrappedFn func()
		err := wrap.After(&wrappedFn, fn, afterFn)
		if err != nil {
			t.Error(err)
		}

		wrappedFn()
		if !afterFnIsCalled {
			t.Error("Third argument must be called")
		}
	})

	t.Run("wrap function with many args and return values", func(t *testing.T) {
		afterFnIsCalled, fnIsCalled := false, false
		fn := func(a, b int) (int, int) {
			fnIsCalled = true
			return a + b, a * b
		}
		afterFn := func() {
			if !fnIsCalled {
				t.Error("Second argument must be called before third argument")
			}
			afterFnIsCalled = true
		}

		var wrappedFn func(a, b int) (int, int)
		err := wrap.After(&wrappedFn, fn, afterFn)
		if err != nil {
			t.Error(err)
		}

		ret1, ret2 := wrappedFn(2, 3)
		if !afterFnIsCalled {
			t.Error("Third argument must be called")
		}

		if ret1 != 5 || ret2 != 6 {
			t.Error("returned value is incorrect")
		}
	})

	t.Run("first argument must be function pointer", func(t *testing.T) {
		fn1 := func() {}
		fn2 := func() {}
		fn3 := func() {}
		err := wrap.After(fn1, fn2, fn3)
		if err == nil {
			t.Error("cannot set other than function pointer for first argument")
		}

		intValue := 0
		fn4 := func() {}
		fn5 := func() {}
		err = wrap.After(&intValue, fn4, fn5)
		if err == nil {
			t.Error("cannot set other than function pointer for first argument")
		}
	})

	t.Run("second argument must be function", func(t *testing.T) {
		var dst func()
		fn := func() {}
		err := wrap.After(dst, 1, fn)
		if err == nil {
			t.Error("cannot set other than function for second argument")
		}
	})

	t.Run("value of first argument must be same type as second argument", func(t *testing.T) {
		var dst func(a int32)
		fn := func(a int64) {}
		err := wrap.After(&dst, fn, func() {})
		if err == nil {
			t.Error("value of first argument must be same type as second argument")
		}
	})
}
