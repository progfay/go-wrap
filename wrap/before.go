package wrap

import (
	"fmt"
	"reflect"
)

func Before(dst, src interface{}, beforeFn func()) error {
	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)

	if vDst.Kind() != reflect.Ptr {
		return fmt.Errorf("First argument must be pointer of function")
	}

	vDst = vDst.Elem()
	if vDst.Kind() != reflect.Func {
		return fmt.Errorf("First argument must be pointer of function")
	}

	if vSrc.Kind() != reflect.Func {
		return fmt.Errorf("Second argument must be function")
	}

	if vDst.Type() != vSrc.Type() {
		return fmt.Errorf("Value of first argument must be same type function as second argument")
	}

	f := func(values []reflect.Value) []reflect.Value {
		beforeFn()
		return vSrc.Call(values)
	}

	vDst.Set(reflect.MakeFunc(vDst.Type(), f))

	return nil
}
